package shardedengine

type ShardedEngine struct {
	engines      []BaseEngine
	shardsAmount int
}

//go:generate mockery --name BaseEngine
type BaseEngine interface {
	Get(key string) (*string, error)
	Set(key, value string)
	Del(key string)
}

func NewShardedEngine(logSize int, newEngine func() BaseEngine) *ShardedEngine {
	shardsAmount := 1 << logSize
	engines := make([]BaseEngine, 0, shardsAmount)
	for range shardsAmount {
		engines = append(engines, newEngine())
	}

	return &ShardedEngine{
		engines:      engines,
		shardsAmount: shardsAmount,
	}
}

func (e *ShardedEngine) Get(key string) (*string, error) {
	hash := e.getHash(key) & uint32(e.shardsAmount-1)
	return e.engines[hash].Get(key)
}

func (e *ShardedEngine) Set(key, value string) {
	hash := e.getHash(key) & uint32(e.shardsAmount-1)
	e.engines[hash].Set(key, value)
}

func (e *ShardedEngine) Del(key string) {
	hash := e.getHash(key) & uint32(e.shardsAmount-1)
	e.engines[hash].Del(key)
}

func (e *ShardedEngine) getHash(key string) uint32 {
	var hash uint32 = 5381

	for _, v := range key {
		hash = ((hash << 5) + hash) + uint32(v)
	}

	return hash
}
