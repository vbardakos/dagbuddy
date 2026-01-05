package buddy

import (
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vbardakos/dagbuddy/lsp"
	"github.com/vbardakos/dagbuddy/rpc"
)

type BuddyManager struct {
	Ctx    *rpc.SessionManager
	Log    *log.Logger
	Shared *SharedState
}

type SharedState struct {
	Client *lsp.InitializeClientParams
	Server *lsp.InitializeServerResult
}

type FileState struct {
	Text string
	Lang string

	version atomic.Int32
	mu      sync.RWMutex
}

func (s SharedState) IsOk() bool {
	return s.Client != nil && s.Server != nil
}

func (f *FileState) Open(txt, lang string, next int32) {
	f.version.Swap(next)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.Text = txt
	f.Lang = lang
}

func (f *FileState) Le(version int32) bool {
	cur := f.version.Load()
	return cur <= version
}

func (f *FileState) CommitChanges(cs []lsp.TextDocumentChanges) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var lines []string
	var lineLen int
	for _, c := range cs {
		if c.Length == 0 && c.Range == nil {
			f.Text = c.Text
		}

		rem, lines := strings.Split(c.Text, "\n"), strings.Split(f.Text, "\n")
		for col := c.Range.Start.Col; col < c.Range.End.Col; col++ {

		}
	}
}
