package engine

import "errors"

type Engine struct {
	shards []*Shard
	size   int
}

func NewEngine(logSize int) (*Engine, error) {
	if logSize < 0 {
		return nil, errors.New("log shards amount must be non-negative")
	}

	size := 1 << logSize
	shards := make([]*Shard, 0, size)
	for range size {
		shards = append(shards, NewShard())
	}

	return &Engine{
		shards: shards,
		size:   size,
	}, nil
}

func (e *Engine) Get(key string) (string, error) {
	hash := e.getHash(key) & uint32(e.size-1)
	return e.shards[hash].Get(key)
}

func (e *Engine) Set(key, value string) {
	hash := e.getHash(key) & uint32(e.size-1)
	e.shards[hash].Set(key, value)
}

func (e *Engine) Del(key string) {
	hash := e.getHash(key) & uint32(e.size-1)
	e.shards[hash].Del(key)
}

func (e *Engine) getHash(key string) uint32 {
	var hash uint32 = 5381

	for _, v := range key {
		hash = ((hash << 5) + hash) + uint32(v)
	}

	return hash
}
