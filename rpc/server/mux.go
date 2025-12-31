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

func (mx *mux) setSessionBuilder(ctx context.Context) context.Context {
	if mx.opts.SessionBuilder == nil {
		return ctx
	}
	ns := mx.opts.SessionBuilder
	mx.opts.SessionBuilder = nil

	sm := newSessionManager(ns)
	return sm.attach(ctx)
}

func (mx *mux) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	ctx = mx.setSessionBuilder(ctx)
	defer cancel()

	in := make(chan p.RPCMessage, mx.opts.InSize)
	out := make(chan p.RPCMessage, mx.opts.OutSize)
	errs := make(chan error, 1)

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
			return ctx.Err()

		case err := <-errs:
			mx.opts.OnError(err)
			if mx.opts.Shutdown == FailFast {
				return err
			}

		case msg, ok := <-in:
			if !ok {
				return nil
			}

			dispatcher.dispatch(ctx, msg)
		}
	}
}

func (mx *mux) SingleHandshake(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	ctx = mx.setSessionBuilder(ctx)
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
