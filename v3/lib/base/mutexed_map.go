package base

import (
	// "log"
	"sync"
	"sync/atomic"
)

type MutexedMap[T any] struct {
	sync.RWMutex
	M     map[string]T `json:"map"`
	count atomic.Int64
}

func NewMutexedMap[T any]() *MutexedMap[T] {

	return &MutexedMap[T]{
		M:     make(map[string]T),
		count: atomic.Int64{},
	}
}

func (mm *MutexedMap[T]) Set(key string, value T) {
	// log.Println("MutexedMap[T] set", key)
	mm.Lock()
	mm.M[key] = value
	mm.Unlock()
	mm.count.Add(1)
	// log.Println("MutexedMap[T] unlocked set", key)
}

func (mm *MutexedMap[T]) Get(key string) (t T, ok bool) {
	// log.Println("MutexedMap[T] ", key)
	mm.RLock()
	t, ok = mm.M[key]
	mm.RUnlock()
	// log.Println("MutexedMap[T] unlocked", key)
	return
}

func (mm *MutexedMap[T]) Clear() {
	mm.Lock()
	clear(mm.M)
	mm.count.Store(0)
	mm.Unlock()

}
func (mm *MutexedMap[T]) Keys() (keys []string) {
	mm.RLock()
	mm.count.Store(0)
	for key := range mm.M {
		keys = append(keys, key)
		mm.count.Add(1)
	}
	mm.RUnlock()
	return
}
func (mm *MutexedMap[T]) Delete(key string) (t T, ok bool) {
	mm.Lock()
	t, ok = mm.M[key]
	if ok {
		delete(mm.M, key)
		mm.count.Add(-1)
	}
	mm.Unlock()
	return
}

func (mm *MutexedMap[T]) Copy() (new_m map[string]T, count int) {
	new_m = make(map[string]T)
	mm.RLock()
	count = 0
	for key := range mm.M {
		new_m[key] = mm.M[key]
		count += 1
	}
	mm.RUnlock()
	return
}
