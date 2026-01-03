package lsp

import (
	"context"
	"log"

	"github.com/vbardakos/dagbuddy/rpc"
)

type LspHandler struct {
	Log    *log.Logger
	server ServerCaps
	client InitializeParams
}

func (h LspHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	h.Log.Println("Enters handling")
	msg, err := rpc.DecodeMessage(data)

	if _, ok := msg.(rpc.NotificationMessage); ok {
		return nil, nil
	}

	if r, ok := msg.(rpc.RequestMessage); ok {
		resp, err := dispatch(ctx, &h, r)
		if err != nil {
			return nil, err
		}
		if msg, ok := resp.(rpc.RPCMessage); ok {
			return rpc.EncodeMessage(msg)
		}
		return nil, rpc.InternalError
	}

	if err != nil {
		h.Log.Printf("Decoding Err: %s\n", err)
		return []byte{}, err
	}
	h.Log.Printf("Decoded Message: %+v\n", msg)
	return data, err
}

type ServerCaps struct{}
type ClientCaps struct{}
