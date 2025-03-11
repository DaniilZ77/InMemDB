package engine

import (
	"sync"
)

type Engine struct {
	mu   sync.Mutex
	data map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) Get(key string) (*string, error) {
	e.mu.Lock()
	value, ok := e.data[key]
	e.mu.Unlock()

	if !ok {
		return nil, ErrKeyNotFound
	}

	return &value, nil
}

func (e *Engine) Set(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.data[key] = value
}

func (e *Engine) Del(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.data, key)
}
