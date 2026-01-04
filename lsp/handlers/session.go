package handlers

import (
	"log"

	s "github.com/vbardakos/dagbuddy/lsp/schemas"
	"github.com/vbardakos/dagbuddy/rpc"
)

var GetLspSession = rpc.GetUserSession[Session]

type Session struct {
	Client *s.InitializeParams
	Server *s.InitializeResult
	Log    *log.Logger
}

func NewSession(l *log.Logger) *Session {
	return &Session{Log: l}
}
