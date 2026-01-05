package lsp

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/vbardakos/dagbuddy/rpc"
)

type LspCodec struct{}

func (c *LspCodec) Read(_ context.Context, r *bufio.Reader) (rpc.Message, error) {
	var msgLen int

	// note :: read header; head & msg split by double \r\n
	for {
		header, err := r.ReadString('\n')

		if err != nil {
			return nil, err
		}

		if header == "\r\n" {
			break
		}

		fmt.Sscanf(header, "Content-Length: %d", &msgLen)
	}

	data := make([]byte, msgLen)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	return rpc.Decode(data)
}

// Dispatcher takes care of NoResponse
func (c *LspCodec) Write(ctx context.Context, w *bufio.Writer, msg rpc.Message) error {
	data, err := rpc.Encode(msg)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(data))
	_, err = w.Write(data)
	return err
}
