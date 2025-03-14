package shardedengine

import "hash/fnv"

type ShardedEngine struct {
	engines []BaseEngine
}

type BaseEngine interface {
	Get(key string) (*string, error)
	Set(key, value string)
	Del(key string)
}

func NewShardedEngine(size int, newEngine func() BaseEngine) *ShardedEngine {
	engines := make([]BaseEngine, 0, size)
	for range size {
		engines = append(engines, newEngine())
	}

	return &ShardedEngine{
		engines: engines,
	}
}

func (e *ShardedEngine) Get(key string) (*string, error) {
	h := fnv.New32a()
	h.Write([]byte(key))

	return e.engines[int(h.Sum32())%len(e.engines)].Get(key)
}

func (e *ShardedEngine) Set(key, value string) {
	h := fnv.New32a()
	h.Write([]byte(key))

	e.engines[int(h.Sum32())%len(e.engines)].Set(key, value)
}

func (e *ShardedEngine) Del(key string) {
	h := fnv.New32a()
	h.Write([]byte(key))

	e.engines[int(h.Sum32())%len(e.engines)].Del(key)
}
