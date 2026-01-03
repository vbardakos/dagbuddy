package server

import (
	s "github.com/vbardakos/dagbuddy/rpc/session"
)

type ShutdownMode int
type FlushMode int

const (
	Drain ShutdownMode = iota
	FailFast

	FlushPerMessage FlushMode = iota
	FlushOnClose
)

type OptionFunc func(*Options)

type Options struct {
	InSize         int
	OutSize        int
	Shutdown       ShutdownMode
	Flush          FlushMode
	OnError        func(error)
	SessionManager s.SessionManager
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

func WithSessionManager(sm func() s.SessionManager) OptionFunc {
	return func(o *Options) {
		o.SessionManager = sm()
	}
}

func WithSession(o *Options) {
	o.SessionManager = s.NewSingleSessionManager()
}

func WithUserSession[ST any](fn func() *ST) OptionFunc {
	return func(o *Options) {
		o.SessionManager = s.NewGenPhantomManager(fn)
	}
}

func WithMultiSessionManager(o *Options) {
	o.SessionManager = s.NewMultiSessionManager()
}

func WithUserMultiSessionManager[ID comparable, ST any](fn func(ID) *ST) OptionFunc {
	return func(o *Options) {
		o.SessionManager = s.NewGenMapManager(fn)
	}
}

func DrainErrors(o *Options) {
	o.Shutdown = Drain
}

func ExitOnError(o *Options) {
	o.Shutdown = FailFast
}

func MaxInSize(n int) OptionFunc {
	return func(o *Options) {
		o.InSize = n
	}
}

func MaxOutSize(n int) OptionFunc {
	return func(o *Options) {
		o.OutSize = n
	}
}

func OnError(fn func(error)) OptionFunc {
	return func(o *Options) {
		o.OnError = fn
	}
}

func WithoutFlush(o *Options) {
	o.Flush = FlushOnClose
}
