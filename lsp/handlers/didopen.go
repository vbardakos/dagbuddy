package handlers

import (
	"context"
	"encoding/json"

	s "github.com/vbardakos/dagbuddy/lsp/schemas"
	"github.com/vbardakos/dagbuddy/rpc"
)

type TextDocDidOpen struct{}

func (*TextDocDidOpen) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	req, ok := msg.(rpc.Notification)
	if !ok {
		return nil, rpc.InternalError
	}

	var didOpen s.DidOpenTextDocument
	if err := json.Unmarshal(req.Params, &didOpen); err != nil {
		return nil, err
	}

	sess := GetLspSession(ctx)
	sess.Log.Printf("didopen params: %+v", didOpen)

	return rpc.NoResponse, nil
}
