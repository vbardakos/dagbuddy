package lsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
)

type lspCodec struct {
	Log *log.Logger
}

func NewLspCodec(l *log.Logger) *lspCodec {
	return &lspCodec{Log: l}
}

func (c lspCodec) Read(_ context.Context, r *bufio.Reader) ([]byte, error) {
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
	c.Log.Printf("Messsage Length: %d\n", msgLen)
	msg := make([]byte, msgLen)
	_, err := io.ReadFull(r, msg)
	c.Log.Printf("Read Messsage: %s. Err?: %s\n", msg, err)
	return msg, err
}

func (c lspCodec) Write(ctx context.Context, w *bufio.Writer, msg []byte) error {
	if msg == nil {
		return nil
	}

	c.Log.Printf("Write Message: %s", msg)
	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(msg))
	_, err := w.Write(msg)
	return err
}
