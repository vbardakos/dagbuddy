package server

import "sync"

func NewDefaultSession(id SessionID) Session {
	return &defaultSession{ID: id, Values: map[string]any{}}
}

type defaultSession struct {
	ID     SessionID
	Values map[string]any
	Mu     sync.RWMutex
}

func (*defaultSession) session() {}
