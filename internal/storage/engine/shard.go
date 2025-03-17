package engine

import (
	"sync"
)

type Shard struct {
	data sync.Map
}

func NewShard() *Shard {
	return &Shard{}
}

func (e *Shard) Get(key string) (string, error) {
	value, ok := e.data.Load(key)
	if !ok {
		return "", ErrKeyNotFound
	}

	return value.(string), nil
}

func (e *Shard) Set(key, value string) {
	e.data.Store(key, value)
}

func (e *Shard) Del(key string) {
	e.data.Delete(key)
}
