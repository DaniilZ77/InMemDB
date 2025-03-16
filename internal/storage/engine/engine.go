package engine

import (
	"sync"
)

type Engine struct {
	data sync.Map
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Get(key string) (string, error) {
	value, ok := e.data.Load(key)
	if !ok {
		return "", ErrKeyNotFound
	}

	return value.(string), nil
}

func (e *Engine) Set(key, value string) {
	e.data.Store(key, value)
}

func (e *Engine) Del(key string) {
	e.data.Delete(key)
}
