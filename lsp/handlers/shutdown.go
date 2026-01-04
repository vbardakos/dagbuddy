package handlers

import (
	"context"

	"github.com/vbardakos/dagbuddy/rpc"
)

type ShutdownHandler struct{}
type ExitNotifier struct{}

func (*ShutdownHandler) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	delete(requests, "shutdown")

	if r, ok := msg.(rpc.Request); ok {
		return rpc.NewResponse(r.ID, make(map[string]any, 0), nil)
	}
	return nil, rpc.InternalError
}

func (*ExitNotifier) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	delete(notifications, "exit")
	rpc.Done(ctx)
	return rpc.NoResponse, nil
}
