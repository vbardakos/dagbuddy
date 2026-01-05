package buddy

import (
	"encoding/json"
	"io"
	"log"

	"github.com/vbardakos/dagbuddy/rpc"
)

type Selector struct {
	log *log.Logger
	rs  map[string]rpc.MessageHandler
	ns  map[string]rpc.MessageHandler
	// resp map[string]rpc.MessageHandler // note :: No use case yet
}

func NewSelector(rs, ns map[string]rpc.MessageHandler, l *log.Logger) rpc.Selector {
	if l == nil {
		l = log.New(io.Discard, "", 0)
	}
	s := &Selector{log: l, rs: rs, ns: ns}
	return s.Select
}

func (s *Selector) Select(msg rpc.Message) rpc.MessageHandler {
	if r, ok := msg.(rpc.Request); ok {
		s.log.Printf("selects request method: %s\n", r.Method)
		return s.rs[r.Method]
	}

	if n, ok := msg.(rpc.Notification); ok {
		s.log.Printf("selects notification method: %s\n", n.Method)
		return s.ns[n.Method]
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
