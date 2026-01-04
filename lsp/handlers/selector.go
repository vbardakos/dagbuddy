package handlers

import (
	"encoding/json"
	"log"

	"github.com/vbardakos/dagbuddy/rpc"
)

var requests = map[string]rpc.MessageHandler{
	"initialize": &Initializer{},
	"shutdown":   &ShutdownHandler{},
}

var notifications = map[string]rpc.MessageHandler{
	"initialized":            &InitializeNotifier{},
	"textDocument/didOpen":   &TextDocDidOpen{},
	"textDocument/didClose":  &TextDocDidClose{},
	"textDocument/didChange": &TextDocDidChange{},
	"exit":                   &ExitNotifier{},
}

type Selector struct {
	log *log.Logger
}

func NewSelector(l *log.Logger) rpc.Selector {
	s := &Selector{log: l}
	return s.selectHandler
}

func (s *Selector) selectHandler(msg rpc.Message) rpc.MessageHandler {
	if r, ok := msg.(rpc.Request); ok {
		s.log.Printf("selects request method: %s\n", r.Method)
		return requests[r.Method]
	}

	if n, ok := msg.(rpc.Notification); ok {
		s.log.Printf("selects notification method: %s\n", n.Method)
		return notifications[n.Method]
	}

	if r, ok := msg.(rpc.Response); ok {
		if r.Error != nil {
			s.log.Printf("selects response error: %+v\n", r.Error)
			return nil
		}

		var res map[string]any
		if err := json.Unmarshal(r.Result, &res); err != nil {
			s.log.Printf("selects response: unmarshal err\n")
			return nil
		}

		s.log.Printf("selects response result: %s\n", res)
	}

	return nil
}
