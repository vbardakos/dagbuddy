package server

import (
	"bufio"
	"context"
	"io"

	p "github.com/vbardakos/dagbuddy/rpc/protocol"
)

type mux struct {
	r     io.Reader
	w     io.Writer
	codec Codec
	opts  Options
	selFn Selector
}

func NewMultiplexer(r io.Reader, w io.Writer, c Codec, s Selector, opts ...OptionFunc) *mux {
	mx := &mux{r: r, w: w, codec: c, opts: defaultOpts(), selFn: s}
	for _, fn := range opts {
		fn(&mx.opts)
	}
	return mx
}

func (mx *mux) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	ctx = mx.attachSessionManager(ctx)
	defer cancel()

	var pendingDone bool
	ctx, doneNow, doneAfter := mx.setupDoneCapabilities(ctx)

	in, out, errs := mx.setupIOChans()
	go mx.readLoop(ctx, in, errs)
	go mx.writeLoop(ctx, out, errs)

	dispatcher := &dispatcher{
		selectFunc: mx.selFn,
		onError:    mx.opts.OnError,
		out:        out,
	}

	for {
		select {

		case <-ctx.Done():
			return nil

		// user immediate shutdown
		case <-doneNow:
			cancel()
			return nil

		// user shutdown after dispatch
		case <-doneAfter:
			pendingDone = true

		// Handle error handling
		case err := <-errs:
			mx.opts.OnError(err)
			if mx.opts.Shutdown == FailFast {
				return err
			}

		// happy path
		case msg, ok := <-in:
			if !ok {
				return nil
			}

			dispatcher.dispatch(ctx, msg)

			if pendingDone {
				cancel()
				return nil
			}
		}
	}
}

// TODO :: use doneAfter for cleaner implementation
func (mx *mux) SingleHandshake(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	ctx = mx.attachSessionManager(ctx)
	defer cancel()

	br := bufio.NewReader(mx.r)
	bw := bufio.NewWriter(mx.w)
	defer bw.Flush()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		recvmsg, err := mx.codec.Read(ctx, br)
		if err != nil {
			mx.opts.OnError(err)
			return err
		}

		handler := mx.selFn(recvmsg)
		if handler == nil {
			continue
		}

		msg, err := handler.Handle(ctx, recvmsg)
		if err != nil {
			mx.opts.OnError(err)
			if mx.opts.Shutdown == FailFast {
				return err
			}
			continue
		}

		if msg.Type() == p.NotificationKind {
			continue
		}

		if err := mx.codec.Write(ctx, bw, msg); err != nil {
			mx.opts.OnError(err)
			return err
		}

		if mx.opts.Flush == FlushPerMessage {
			if err := bw.Flush(); err != nil {
				mx.opts.OnError(err)
				return err
			}
		}

		return nil
	}
}

func (mx *mux) readLoop(ctx context.Context, in chan<- p.RPCMessage, errs chan<- error) {
	defer close(in)

	br := bufio.NewReader(mx.r)
	for {
		msg, err := mx.codec.Read(ctx, br)
		if err != nil {
			errs <- err
			return
		}

		select {
		case <-ctx.Done():
			return

		case in <- msg:
		}
	}
}

func (mx *mux) writeLoop(ctx context.Context, out <-chan p.RPCMessage, errs chan<- error) {
	bw := bufio.NewWriter(mx.w)
	defer bw.Flush()

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-out:
			if !ok {
				return
			}

			if err := mx.codec.Write(ctx, bw, msg); err != nil {
				errs <- err
				return
			}

			if mx.opts.Flush == FlushPerMessage {
				if err := bw.Flush(); err != nil {
					errs <- err
					return
				}
			}
		}
	}
}

func (mx *mux) attachSessionManager(ctx context.Context) context.Context {
	if mx.opts.SessionManager == nil {
		return ctx
	}
	sm := mx.opts.SessionManager
	mx.opts.SessionManager = nil

	return sm.Attach(ctx)
}

func (mx *mux) setupIOChans() (in chan p.RPCMessage, out chan p.RPCMessage, errs chan error) {
	in = make(chan p.RPCMessage, mx.opts.InSize)
	out = make(chan p.RPCMessage, mx.opts.OutSize)
	errs = make(chan error, 1)
	return
}

func (mx *mux) setupDoneCapabilities(ctxin context.Context) (ctx context.Context, done chan struct{}, doneAfter chan struct{}) {
	done = make(chan struct{}, 1)
	ctx = attachDoneNowContextValue(ctxin, done)

	doneAfter = make(chan struct{}, 1)
	ctx = attachDoneAfterContextValue(ctx, doneAfter)
	return
}
