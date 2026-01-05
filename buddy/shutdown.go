package buddy

import (
	"context"

	"github.com/vbardakos/dagbuddy/rpc"
)

type ShutdownRequest struct{}
type ExitNotification struct{}

func (*ShutdownRequest) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	if r, ok := msg.(rpc.Request); ok {
		return rpc.NewResponse(r.ID, make(map[string]any, 0), nil)
	}
	return nil, rpc.InternalError
}

func (*ExitNotification) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	rpc.Done(ctx)
	return rpc.NoResponse, nil
}
