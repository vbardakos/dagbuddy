package session

import (
	"context"

	i "github.com/vbardakos/dagbuddy/rpc/session/internal"
)

type mapManagerID struct{}

type MapManager[ID comparable, ST any] struct {
	ss    *i.SafeMap[ID, *ST]
	newFn func(ID) *ST
}

func NewGenMapManager[ID comparable, ST any](fn func(ID) *ST) *MapManager[ID, ST] {
	return &MapManager[ID, ST]{ss: i.NewMap[ID, *ST](), newFn: fn}
}

func GetGenMapFromContext[ID comparable, T any](ctx context.Context) *MapManager[ID, T] {
	return ctx.Value(mapManagerID{}).(*MapManager[ID, T])
}

func (sm *MapManager[ID, ST]) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, mapManagerID{}, sm)
}

func (sm *MapManager[ID, ST]) Get(id ID) (*ST, bool) {
	return sm.ss.Get(id)
}

func (sm *MapManager[ID, ST]) Del(id ID) {
	sm.ss.Del(id)
}

func (sm *MapManager[ID, ST]) Len(id ID) int {
	return sm.ss.Len()
}

func (sm *MapManager[ID, ST]) New(id ID) *ST {
	sess := sm.newFn(id)
	sm.ss.Set(id, sess)
	return sess
}

func (sm *MapManager[ID, ST]) Clear() {
	sm.ss = i.NewMap[ID, *ST]()
}

func (sm *MapManager[ID, ST]) Apply(fn func(id ID, sess *ST)) {
	sm.ss.Mu.Lock()
	defer sm.ss.Mu.Unlock()
	for id, s := range sm.ss.Data {
		fn(id, s)
	}
}

func (sm *MapManager[ID, ST]) ForEach(fn func(id ID, sess ST)) {
	sm.ss.Mu.RLock()
	defer sm.ss.Mu.RUnlock()
	for id, s := range sm.ss.Data {
		fn(id, *s)
	}
}
