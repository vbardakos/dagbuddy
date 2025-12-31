package server

import (
	"bufio"
	"context"

	p "github.com/vbardakos/dagbuddy/rpc/protocol"
)

type Codec interface {
	Read(ctx context.Context, r *bufio.Reader) (p.RPCMessage, error)
	Write(ctx context.Context, w *bufio.Writer, msg p.RPCMessage) error
}
