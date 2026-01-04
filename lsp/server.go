package lsp

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	h "github.com/vbardakos/dagbuddy/lsp/handlers"
	"github.com/vbardakos/dagbuddy/rpc"
)

func RunLspServer(l *log.Logger, o ...rpc.UserOption) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	r, w := stdIO()
	o = append(o, rpc.WithUserSingleSession(wrappedSession(l)))

	s := rpc.NewServer(r, w, &LspCodec{Log: l}, h.NewSelector(l), o...)
	return s.Run(ctx)
}

func stdIO() (io.Reader, io.Writer) {
	return os.Stdin, os.Stdout
}

func wrappedSession(l *log.Logger) func() *h.Session {
	return func() *h.Session { return h.NewSession(l) }
}
