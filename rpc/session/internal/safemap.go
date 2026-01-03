package internal

import "sync"

type SafeMap[KT comparable, VT any] struct {
	Data map[KT]VT
	Mu   sync.RWMutex
}

func NewMap[KT comparable, VT any]() *SafeMap[KT, VT] {
	return &SafeMap[KT, VT]{Data: map[KT]VT{}}
}

func (m *SafeMap[KT, VT]) Get(k KT) (VT, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	v, ok := m.Data[k]
	return v, ok
}

func (m *SafeMap[KT, VT]) Set(k KT, v VT) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Data[k] = v
}

func (m *SafeMap[KT, VT]) Del(k KT) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	delete(m.Data, k)
}

func (m *SafeMap[KT, VT]) Len() int {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	return len(m.Data)
}
