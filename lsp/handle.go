package lsp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/vbardakos/dagbuddy/rpc"
)

type LspHandler struct {
	Log *log.Logger
}

func (h LspHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	h.Log.Println("Enters handling")
	msg, err := rpc.DecodeMessage(data)
	if r, ok := msg.(rpc.RequestMessage); ok {
		if r.Method == "initialize" {
			h.Log.Println("Found initialize Request")
			var params InitializeParams
			err := json.Unmarshal(r.Params, &params)
			if err != nil {
				h.Log.Printf("Failed to marshal params. %s\n", err)
			}
			h.Log.Println("Successfully marshalled params")
			h.Log.Printf("%+v", params)
		}
	}
	if err != nil {
		h.Log.Printf("Decoding Err: %s\n", err)
		return []byte{}, err
	}
	h.Log.Printf("Decoded Message: %+v\n", msg)
	return data, err
}
