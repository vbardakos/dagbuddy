package analysis

// import (
// 	"sync"
// 	"sync/atomic"
//
// 	"github.com/vbardakos/dagbuddy/rpc"
// )
//
// type MethodHandlerMap = map[string]rpc.MessageHandler
//
// type StagedHandlerMap struct {
// 	stage     atomic.Uint32
// 	stages    []*defaultMap
// 	len       atomic.Uint32
// 	OutBounds rpc.MessageHandler
// }
//
// func NewStagedHandlerMap(outbounds rpc.MessageHandler) *StagedHandlerMap {
// 	return &StagedHandlerMap{stages: []*defaultMap{}, OutBounds: outbounds}
// }
//
// func (hm *StagedHandlerMap) Append(m MethodHandlerMap, def rpc.MessageHandler) {
// 	hm.len.Add(1)
// 	hm.stages = append(hm.stages, newDefMap(m, def))
// }
//
// func (hm *StagedHandlerMap) Get(md string) rpc.MessageHandler {
// 	if hm.len. {
// 		return hm.OutBounds
// 	}
// }
//
// func (hm *StagedHandlerMap) Get(method string) rpc.MessageHandler {
// 	hm.mu.RLock()
// 	defer hm.mu.RUnlock()
// 	return hm.head.get(method)
// }
//
// func (hm *StagedHandlerMap) Next() bool {
// 	hm.mu.Lock()
// 	defer hm.mu.Unlock()
// 	if inner, ok := hm.tail[0]; ok {
// 		hm.head = inner
// 		hm.tail = hm.tail[1 : len(hm.tail)-1]
// 		return true
// 	}
// 	return false
// }
//
// func newDefMap(vs MethodHandlerMap, def rpc.MessageHandler) *defaultMap {
// 	return &defaultMap{
// 		values: vs,
// 		def:    def,
// 	}
// }
//
// type defaultMap struct {
// 	values map[string]rpc.MessageHandler
// 	def    rpc.MessageHandler
// }
//
// func (m defaultMap) get(key string) rpc.MessageHandler {
// 	if v, ok := m.values[key]; ok {
// 		return v
// 	}
// 	return m.def
// }
