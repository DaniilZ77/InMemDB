package storage

import (
	"math/rand/v2"
	"runtime"
	"sync"
	"testing"

	"github.com/DaniilZ77/InMemDB/internal/storage/baseengine"
	"github.com/DaniilZ77/InMemDB/internal/storage/shardedengine"
)

const (
	alphabet     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shardsAmount = 10
	strsAmount   = 1000
	strLen       = 50
)

func genString(maxSize int) string {
	n := rand.IntN(maxSize) + 1

	var str []byte
	for range n {
		str = append(str, alphabet[rand.IntN(len(alphabet))])
	}

	return string(str)
}

func genStrings(amount, size int) []string {
	strs := make([]string, 0, amount)
	for range amount {
		strs = append(strs, genString(size))
	}

	return strs
}

type engineSetter interface {
	Set(key, value string)
}

func fillEngine(engine engineSetter, strs []string) {
	for i, str := range strs {
		engine.Set(str, strs[i])
	}
}

func BenchmarkEngine_Get(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := baseengine.NewEngine()

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Get(strs[i%strsAmount])
			}
		}()
	}

	wg.Wait()
}

func BenchmarkShardedEngine_Get(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := shardedengine.NewShardedEngine(shardsAmount, func() shardedengine.BaseEngine {
		return baseengine.NewEngine()
	})

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Get(strs[i%strsAmount])
			}
		}()
	}

	wg.Wait()
}

func BenchmarkEngine_Set(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := baseengine.NewEngine()

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Set(strs[i%strsAmount], strs[(i+1)%strsAmount])
			}
		}()
	}

	wg.Wait()
}

func BenchmarkShardedEngine_Set(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := shardedengine.NewShardedEngine(shardsAmount, func() shardedengine.BaseEngine {
		return baseengine.NewEngine()
	})

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Set(strs[i%strsAmount], strs[(i+1)%strsAmount])
			}
		}()
	}

	wg.Wait()
}

func BenchmarkEngine_Del(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := baseengine.NewEngine()

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Del(strs[i%strsAmount])
			}
		}()
	}

	wg.Wait()
}

func BenchmarkShardedEngine_Del(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := shardedengine.NewShardedEngine(shardsAmount, func() shardedengine.BaseEngine {
		return baseengine.NewEngine()
	})

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Del(strs[i%strsAmount])
			}
		}()
	}

	wg.Wait()
}
