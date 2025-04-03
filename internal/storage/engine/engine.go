package engine

import (
	"errors"
	"hash/fnv"
)

type Engine struct {
	shards []*Shard
}

func NewEngine(shardsNumber int) (*Engine, error) {
	if shardsNumber < 1 {
		return nil, errors.New("shards number must be positive")
	}

	shards := make([]*Shard, 0, shardsNumber)
	for range shardsNumber {
		shards = append(shards, NewShard())
	}

	return &Engine{
		shards: shards,
	}, nil
}

func (e *Engine) Get(key string) (string, bool) {
	return e.shards[e.getHash(key)].Get(key)
}

func (e *Engine) Set(key, value string) {
	e.shards[e.getHash(key)].Set(key, value)
}

func (e *Engine) Del(key string) {
	e.shards[e.getHash(key)].Del(key)
}

func (e *Engine) getHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32() % uint32(len(e.shards))
}
