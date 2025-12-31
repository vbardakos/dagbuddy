package server

import (
	"context"
	"sync"
)

type SessionID string
type managerID struct{}

type Session interface {
	session() // ensures ptr
}

type sessionManager struct {
	ss map[SessionID]Session
	mu sync.RWMutex
	ns func(id SessionID) Session
}

func SessionManagerFromContext(ctx context.Context) *sessionManager {
	sm, _ := ctx.Value(managerID{}).(*sessionManager)
	return sm
}

func (sm *sessionManager) Get(id SessionID) Session {
	sm.mu.RLock()
	if s, ok := sm.ss[id]; ok {
		sm.mu.RUnlock()
		return s
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	s := sm.ns(id)
	sm.ss[id] = s
	return s
}

func (sm *sessionManager) Del(id SessionID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.ss, id)
}

func (sm *sessionManager) Exists(id SessionID) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.ss[id]
	return ok
}

func newSessionManager(ns func(id SessionID) Session) *sessionManager {
	return &sessionManager{ns: ns, ss: make(map[SessionID]Session)}
}

func (sm *sessionManager) attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, managerID{}, sm)
}
