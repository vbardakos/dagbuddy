package handlers

import (
	"context"
	"encoding/json"

	"github.com/vbardakos/dagbuddy/rpc"
)

type TextDocDidChange struct{}

func (*TextDocDidChange) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	req, ok := msg.(rpc.Notification)
	if !ok {
		return nil, rpc.InternalError
	}

	var didChange map[string]any
	if err := json.Unmarshal(req.Params, &didChange); err != nil {
		return nil, err
	}

	sess := GetLspSession(ctx)
	sess.Log.Printf("didChange params: %+v", didChange)

	return rpc.NoResponse, nil
}
