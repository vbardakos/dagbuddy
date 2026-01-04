package handlers

import (
	"context"
	"encoding/json"

	s "github.com/vbardakos/dagbuddy/lsp/schemas"
	"github.com/vbardakos/dagbuddy/rpc"
)

type Initializer struct{}
type InitializeNotifier struct{}

func (*Initializer) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	delete(requests, "initialize")
	rq, ok := msg.(rpc.Request)
	if !ok {
		return nil, rpc.InternalError
	}

	sess := GetLspSession(ctx)

	var initParams s.InitializeParams
	if err := json.Unmarshal(rq.Params, &initParams); err != nil {
		sess.Log.Printf("Parse error on %+v", rq.Params)
		return nil, err
	}

	sess.Client = &initParams
	sess.Server = newDefaultInitResult()

	return rpc.NewResponse(rq.ID, sess.Server, nil)
}

// sanity checks
func (*InitializeNotifier) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	delete(notifications, "initialized")
	s := GetLspSession(ctx)
	if s.Server == nil || s.Client == nil {
		return nil, rpc.InternalError
	}

	s.Log.Println("successfull initialized notification")
	return rpc.NoResponse, nil
}

func newDefaultInitResult() *s.InitializeResult {
	return &s.InitializeResult{
		Caps: &s.ServerCapabilities{
			TextDocSync: s.TextDocSync{
				OpenClose: true,
				Change:    s.Full,
			},
			Completion: &s.CompletionOptions{
				TriggerChars: []string{".", ":", "?"},
				Resolve:      true,
			},
			Hover:      true,
			Definition: true,
		},
		Info: &s.ServerInfo{
			Name:    "yoloLSP",
			Version: "0.0.1-alpha",
		},
	}
}
