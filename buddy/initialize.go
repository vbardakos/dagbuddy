package buddy

import (
	"context"
	"encoding/json"

	"github.com/vbardakos/dagbuddy/lsp"
	"github.com/vbardakos/dagbuddy/rpc"
)

const CAPABILITIES = "__caps"

type InitializeRequest struct{}
type InitializedNotification struct{}

func (*InitializeRequest) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	r, ok := msg.(rpc.Request)
	if !ok {
		return nil, rpc.InternalError
	}

	var initParams lsp.InitializeClientParams
	if err := json.Unmarshal(r.Params, &initParams); err != nil {
		return nil, err
	}

	m := GetStatesManager(ctx)
	s := m.New(CAPABILITIES)
	s.client = &initParams
	s.server = newDefaultInitResult()

	return rpc.NewResponse(r.ID, s.server, nil)
}

// sanity check
func (*InitializedNotification) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	m := GetStatesManager(ctx)
	s, ok := m.Get(CAPABILITIES)
	if !ok || s.server == nil || s.client == nil {
		return nil, rpc.InternalError
	}

	return rpc.NoResponse, nil
}

func newDefaultInitResult() *lsp.InitializeServerResult {
	return &lsp.InitializeServerResult{
		Caps: &lsp.ServerCapabilities{
			TextDocSync: lsp.TextDocSync{
				OpenClose: true,
				Change:    lsp.Full,
			},
			Completion: &lsp.CompletionOptions{
				TriggerChars: []string{".", ":", "?"},
				Resolve:      true,
			},
			Hover:      true,
			Definition: true,
		},
		Info: &lsp.ServerInfo{
			Name:    "yoloLSP",
			Version: "0.0.1-alpha",
		},
	}
}
