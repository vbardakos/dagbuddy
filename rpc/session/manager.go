package session

import (
	"context"

	i "github.com/vbardakos/dagbuddy/rpc/session/internal"
)

type SessionManager interface {
	Attach(ctx context.Context) context.Context
}

func NewSingleSessionManager() *phantomManager[i.SafeMap[string, any]] {
	return NewGenPhantomManager(i.NewMap[string, any])
}

func NewMultiSessionManager() *MapManager[string, i.SafeMap[string, any]] {
	return NewGenMapManager(func(_ string) *i.SafeMap[string, any] {
		return i.NewMap[string, any]()
	})
}
