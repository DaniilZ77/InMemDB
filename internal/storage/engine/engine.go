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

func (e *Engine) Get(version int64, key string) (string, bool) {
	return e.shards[e.getHash(key)].Get(version, key)
}

func (e *Engine) SetMany(version int64, modified map[string]*string) {
	for key, value := range modified {
		e.shards[e.getHash(key)].Set(version, key, value)
	}
}

func (e *Engine) Set(version int64, key string, value *string) {
	e.shards[e.getHash(key)].Set(version, key, value)
}

func (e *Engine) ExistsBetween(version1, version2 int64, key string) bool {
	return e.shards[e.getHash(key)].ExistsBetween(version1, version2, key)
}

func (e *Engine) getHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32() % uint32(len(e.shards))
}
