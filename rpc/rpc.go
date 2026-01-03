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
)

var (
	NewServer = server.NewMultiplexer

	// Options
	WithSingleSession       = server.WithSession
	WithMultiSessionManager = server.WithMultiSessionManager
	WithUserManager         = server.WithSessionManager
	WithErrorDrain          = server.DrainErrors
	WithErrorFail           = server.ExitOnError
	WithInCapacity          = server.MaxInSize
	WithOutCapacity         = server.MaxOutSize
	WithOnError             = server.OnError
	WithoutFlushing         = server.WithoutFlush

	NewTypeSelector    = server.SelectByType
	NewHanlderSelector = server.NewHandlerSelector

	NewResponse     = protocol.NewResponse
	NewRequest      = protocol.NewRequest
	NewNotification = protocol.NewNotification
	InternalError   = protocol.InternalError
	ParseError      = protocol.ParseError

	GetSingleSession       = session.GetPhantomFromContext[any]
	GetMultiSessionManager = session.GetMapFromContext[any]

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

func WithUserSingleSession[ST any](fn func() *ST) UserOption {
	return server.WithUserSession(fn)
}

func WithUserMultiSession[ID comparable, ST any](fn func(ID) *ST) UserOption {
	return server.WithUserMultiSessionManager(fn)
}

func GetUserSession[ST any](ctx context.Context) *ST {
	return session.GetPhantomFromContext[ST](ctx)
}

func GetUserMultiSession[ST any](ctx context.Context) *ST {
	return session.GetMapFromContext[ST](ctx)
}
