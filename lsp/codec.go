package lsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

type LspCodec struct{}

func (LspCodec) Read(_ context.Context, r *bufio.Reader) ([]byte, error) {
	var msgLen int
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

	msg := make([]byte, msgLen)
	_, err := io.ReadFull(r, msg)
	return msg, err
}

func (LspCodec) Write(ctx context.Context, w *bufio.Writer, msg []byte) error {
	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(msg))
	_, err := w.Write(msg)
	return err
}

func NewLspCodec() *LspCodec {
	return &LspCodec{}
}

