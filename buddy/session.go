package buddy

import (
	"sync"

	"github.com/vbardakos/dagbuddy/lsp"
	"github.com/vbardakos/dagbuddy/rpc"
)

var GetStatesManager = rpc.GetUserSessionManager[string, State]
var WithStateManager = rpc.WithUserSessionManager(InitializeState)

type State struct {
	client  *lsp.InitializeClientParams
	server  *lsp.InitializeServerResult
	text    string
	version int
	lang    string
	mu      sync.RWMutex
}

func InitializeState(string) *State {
	return &State{}
}

func (s *State) New(txt string, v int, lang string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.text = txt
	s.version = v
	s.lang = lang
}

func (s *State) UpdateFull(txt string, v int) {
	if s.Version() > v {
		return
	}
	s.New(txt, v, s.lang)
}

func (s *State) Version() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.version
}

func (s *State) Text() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.text
}
