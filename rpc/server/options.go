package server

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
	SessionBuilder func(id SessionID) Session
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

func WithSession(fn func(id SessionID) Session) OptionFunc {
	return func(o *Options) {
		o.SessionBuilder = fn
	}
}

func WithDefaultManager(o *Options) {
	WithSession(NewDefaultSession)(o)
}

func DrainErrors(o *Options) {
	o.Shutdown = Drain
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
