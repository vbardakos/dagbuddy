package server

import (
	"context"

	p "github.com/vbardakos/dagbuddy/rpc/protocol"
)

type Handler interface {
	Handle(ctx context.Context, msg p.RPCMessage) (p.RPCMessage, error)
}

type Selector func(msg p.RPCMessage) Handler

func SelectByType(req Handler, resp Handler, not Handler) Selector {
	return func(msg p.RPCMessage) Handler {
		switch msg.Type() {
		case p.RequestKind:
			return req
		case p.ResponseKind:
			return resp
		case p.NotificationKind:
			return not
		default:
			return nil
		}
	}
}

func NewHandlerSelector(h Handler) Selector {
	return func(msg p.RPCMessage) Handler { return h }
}
