package rpc

import (
	"bufio"
	"context"
	"io"
)

type Codec interface {
	Read(ctx context.Context, r *bufio.Reader) ([]byte, error)
	Write(ctx context.Context, w *bufio.Writer, msg []byte) error
}

type Handler interface {
	Handle(ctx context.Context, msg []byte) ([]byte, error)
}

type ShutdownMode int

const (
	Drain ShutdownMode = iota
	FailFast
)

type FlushMode int

const (
	FlushPerMessage FlushMode = iota
	FlushOnClose
)

type Option func(*Options)

type Options struct {
	InSize   int
	OutSize  int
	Shutdown ShutdownMode
	Flush    FlushMode
	OnError  func(error)
}

type Server struct {
	r     io.Reader
	w     io.Writer
	codec Codec
	opts  Options
}

func NewServer(r io.Reader, w io.Writer, c Codec, opts ...Option) *Server {
	srv := &Server{r: r, w: w, codec: c, opts: defaultOpts()}
	for _, optfn := range opts {
		optfn(&srv.opts)
	}
	return srv
}

func (srv *Server) Read(ctx context.Context, out chan<- []byte, errs chan<- error) {
	defer close(out)

	brd := bufio.NewReader(srv.r)
	for {
		msg, err := srv.codec.Read(ctx, brd)
		if err != nil {
			errs <- err
			return
		}

		select {
		case <-ctx.Done():
			return

		case out <- msg:
		}
	}
}

func (srv *Server) Write(ctx context.Context, in <-chan []byte, errs chan<- error) {
	bw := bufio.NewWriter(srv.w)
	defer bw.Flush()

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-in:
			if !ok {
				return
			}

			if err := srv.codec.Write(ctx, bw, msg); err != nil {
				errs <- err
				return
			}

			if srv.opts.Flush == FlushPerMessage {
				if err := bw.Flush(); err != nil {
					errs <- err
					return
				}
			}
		}
	}
}

func (srv *Server) Run(ctx context.Context, handler Handler) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	in := make(chan []byte, srv.opts.InSize)
	out := make(chan []byte, srv.opts.OutSize)
	errs := make(chan error, 1)

	go srv.Read(ctx, in, errs)
	go srv.Write(ctx, out, errs)

	for {
		select {
		case <-ctx.Done():
			close(out)
			return nil

		case err := <-errs:
			srv.opts.OnError(err)
			if srv.opts.Shutdown == FailFast {
				cancel()
				close(out)
				return err
			}

		case msg, ok := <-in:
			if !ok {
				close(out)
				return nil
			}

			response, err := handler.Handle(ctx, msg)
			if err != nil {
				srv.opts.OnError(err)
				if srv.opts.Shutdown == FailFast {
					cancel()
					close(out)
					return err
				}
			}

			out <- response
		}
	}
}

func (srv *Server) RunOnce(ctx context.Context, handler Handler) error {
	in := make(chan []byte, 1)
	out := make(chan []byte, 1)
	errs := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go srv.Read(ctx, in, errs)
	go srv.Write(ctx, out, errs)

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-errs:
		srv.opts.OnError(err)
		return err

	case msg, ok := <-in:
		if !ok {
			return nil
		}

		resp, err := handler.Handle(ctx, msg)
		if err != nil {
			srv.opts.OnError(err)
			if srv.opts.Shutdown == FailFast {
				return err
			}
		}

		out <- resp
		close(out)
		return nil
	}
}

func defaultOpts() Options {
	return Options{
		InSize:   0,
		OutSize:  0,
		Shutdown: FailFast,
		Flush:    FlushPerMessage,
		OnError:  func(error) {},
	}
}
