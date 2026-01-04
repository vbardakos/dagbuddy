package handlers

import (
	"context"
	"encoding/json"

	"github.com/vbardakos/dagbuddy/rpc"
)

type Completion struct{}

func (*Completion) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	rq, ok := msg.(rpc.Request)
	if !ok {
		return nil, rpc.InternalError
	}

	var params map[string]any
	if err := json.Unmarshal(rq.Params, &params); err != nil {
		return nil, err
	}

	sess := GetLspSession(ctx)
	sess.Log.Printf("completion params: %s\n", params)

	return rpc.NoResponse, nil
}
