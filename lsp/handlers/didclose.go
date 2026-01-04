package handlers

import (
	"context"

	"github.com/vbardakos/dagbuddy/rpc"
)

type TextDocDidClose struct{}

func (*TextDocDidClose) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	return rpc.NoResponse, nil
}
