package rpc

import (
	"context"

	"github.com/vbardakos/dagbuddy/rpc/codec"
	"github.com/vbardakos/dagbuddy/rpc/protocol"
	"github.com/vbardakos/dagbuddy/rpc/server"
	"github.com/vbardakos/dagbuddy/rpc/session"
)

type (
	Message      = protocol.RPCMessage
	Request      = protocol.RequestMessage
	Response     = protocol.ResponseMessage
	Notification = protocol.NotificationMessage
	Error        = protocol.ResponseError

	Selector       = server.Selector
	MessageHandler = server.Handler
	UserOption     = server.OptionFunc
	UserCodec      = server.Codec
	Options        = server.Options

	SessionManager = session.SessionManager
)

var (
	NewServer = server.NewMultiplexer

	// Options
	WithDefaultSession = server.WithSession
	WithDefaultManager = server.WithMultiSessionManager
	WithNewManager     = server.WithSessionManager
	WithErrorDrain     = server.DrainErrors
	WithErrorFail      = server.ExitOnError
	WithInCapacity     = server.MaxInSize
	WithOutCapacity    = server.MaxOutSize
	WithOnError        = server.OnError
	WithoutFlushing    = server.WithoutFlush

	NewTypeSelector    = server.SelectByType
	NewHanlderSelector = server.NewHandlerSelector

	NewResponse     = protocol.NewResponse
	NewRequest      = protocol.NewRequest
	NewNotification = protocol.NewNotification
	NoResponse      = protocol.NoResponseMessage{}
	InternalError   = protocol.InternalError
	ParseError      = protocol.ParseError

	NewDefaultSession = session.NewSingleSessionManager
	NewDefaultManager = session.NewMultiSessionManager
	GetDefaultSession = session.GetPhantomFromContext[any]
	GetDefaultManager = session.GetGenMapFromContext[string, any]

	Encode = codec.EncodeMessage
	Decode = codec.DecodeMessage

	Done      = server.Done
	DoneAfter = server.DoneAfterDispatch
	ExitNow   = server.Done
	ExitAfter = server.DoneAfterDispatch
)

const (
	RequestKind      = protocol.RequestKind
	ResponseKind     = protocol.ResponseKind
	NotificationKind = protocol.NotificationKind
)

func WithUserSession[ST any](fn func() *ST) UserOption {
	return server.WithUserSession(fn)
}

func WithUserSessionManager[ID comparable, ST any](fn func(ID) *ST) UserOption {
	return server.WithUserMultiSessionManager(fn)
}

func GetUserSession[ST any](ctx context.Context) *ST {
	return session.GetPhantomFromContext[ST](ctx)
}

func GetUserSessionManager[ID comparable, ST any](ctx context.Context) *session.MapManager[ID, ST] {
	return session.GetGenMapFromContext[ID, ST](ctx)
}
