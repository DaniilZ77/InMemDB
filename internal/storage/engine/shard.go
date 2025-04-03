package engine

import "sync"

type Shard struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewShard() *Shard {
	return &Shard{
		data: make(map[string]string),
	}
}

func (e *Shard) Get(key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	value, ok := e.data[key]
	return value, ok
}

func (e *Shard) Set(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data[key] = value
}

func (e *Shard) Del(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.data, key)
}
