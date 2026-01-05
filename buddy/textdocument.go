package buddy

import (
	"context"
	"encoding/json"

	"github.com/vbardakos/dagbuddy/lsp"
	"github.com/vbardakos/dagbuddy/rpc"
)

type TextDocDidOpenNotification struct{}
type TextDocDidChangeNotification struct{}
type TextDocDidCloseNotification struct{}

func (*TextDocDidOpenNotification) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	n, ok := msg.(rpc.Notification)
	if !ok {
		return nil, rpc.InternalError
	}

	var didOpen lsp.TextDocumentDidOpenNotificationParams
	if err := json.Unmarshal(n.Params, &didOpen); err != nil {
		return nil, err
	}

	m := GetStatesManager(ctx)
	s, ok := m.Get(didOpen.TextDoc.URI)
	if !ok {
		s = m.New(didOpen.TextDoc.URI)
	}

	s.UpdateFull(didOpen.TextDoc.Text, didOpen.TextDoc.Version)
	return rpc.NoResponse, nil
}

func (*TextDocDidChangeNotification) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	n, ok := msg.(rpc.Notification)
	if !ok {
		return nil, rpc.InternalError
	}

	var didChange lsp.TextDocumentDidChangeNotificationParams
	if err := json.Unmarshal(n.Params, &didChange); err != nil {
		return nil, err
	}

	m := GetStatesManager(ctx)
	if s, ok := m.Get(didChange.TextDoc.URI); ok {
		if s.Version() <= didChange.TextDoc.Version {
			s.mu.Lock()
			defer s.mu.Unlock()
			for _, c := range didChange.Changes {
				s.text = c.Text
			}
		}
		return rpc.NoResponse, nil
	}
	return nil, rpc.InternalError
}

func (*TextDocDidCloseNotification) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	n, ok := msg.(rpc.Notification)
	if !ok {
		return nil, rpc.InternalError
	}

	var didClose lsp.TextDocumentDidCloseNotificationParams
	if err := json.Unmarshal(n.Params, &didClose); err != nil {
		return nil, err
	}

	m := GetStatesManager(ctx)
	m.Del(didClose.TextDoc.URI)

	return rpc.NoResponse, nil
}
