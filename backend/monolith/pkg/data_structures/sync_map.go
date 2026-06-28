package datastructures

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	m sync.Map
	// p sync.Pool
}

func (tm *SyncMap[K, V]) Load(key K) (V, bool) {
	v, ok := tm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}

	return v.(V), true
}

func (tm *SyncMap[K, V]) Store(key K, value V) {
	tm.m.Store(key, value)
}

func (tm *SyncMap[K, V]) Delete(key K) {
	tm.m.Delete(key)
}

func (tm *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := tm.m.LoadOrStore(key, value)
	return actual.(V), loaded
}
