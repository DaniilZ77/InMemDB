package engine

import (
	"sync"
)

type Shard struct {
	mu   sync.RWMutex
	data *SkipList
}

func NewShard() *Shard {
	return &Shard{
		data: NewSkipList(nil),
	}
}

func (e *Shard) Get(version int64, key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	value, ok := e.data.Find(VersionedKey{
		key:     key,
		version: version,
	})
	return value, ok
}

func (e *Shard) Set(version int64, key string, value *string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data.Insert(VersionedKey{
		key:     key,
		version: version,
	}, value)
}

func (e *Shard) ExistsBetween(version1, version2 int64, key string) bool {
	lower := VersionedKey{
		key:     key,
		version: version1,
	}
	upper := VersionedKey{
		key:     key,
		version: version2,
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.data.ExistsBetween(lower, upper)
}
