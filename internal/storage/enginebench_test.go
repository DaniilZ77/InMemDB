package storage

import (
	"math/rand/v2"
	"runtime"
	"sync"
	"testing"

	"github.com/DaniilZ77/InMemDB/internal/storage/baseengine"
	"github.com/DaniilZ77/InMemDB/internal/storage/shardedengine"
)

// logShardsAmount = 3
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	23879210	        44.20 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Get-11    	12668874	        87.97 ns/op	       0 B/op	       0 allocs/op
// BenchmarkEngine_Set-11           	  492286	      2247 ns/op	     533 B/op	      33 allocs/op
// BenchmarkShardedEngine_Set-11    	 1305135	       953.3 ns/op	     530 B/op	      33 allocs/op
// BenchmarkEngine_Del-11           	24654747	        40.57 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	15178352	        73.99 ns/op	       0 B/op	       0 allocs/op

// logShardsAmount = 4
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	22526079	        44.43 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Get-11    	15148756	        71.00 ns/op	       0 B/op	       0 allocs/op
// BenchmarkEngine_Set-11           	  499623	      2160 ns/op	     533 B/op	      33 allocs/op
// BenchmarkShardedEngine_Set-11    	 1635234	       731.4 ns/op	     529 B/op	      33 allocs/op
// BenchmarkEngine_Del-11           	28006804	        41.03 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	16315332	        80.70 ns/op	       0 B/op	       0 allocs/op

// logShardsAmount = 5
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	22017288	        46.28 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Get-11    	14888175	        75.58 ns/op	       0 B/op	       0 allocs/op
// BenchmarkEngine_Set-11           	  497182	      2060 ns/op	     533 B/op	      33 allocs/op
// BenchmarkShardedEngine_Set-11    	 1887160	       624.0 ns/op	     529 B/op	      33 allocs/op
// BenchmarkEngine_Del-11           	28990063	        43.21 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	14786041	        71.08 ns/op	       0 B/op	       0 allocs/op

// logShardsAmount = 6
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	21745136	        47.68 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Get-11    	14265547	        73.09 ns/op	       0 B/op	       0 allocs/op
// BenchmarkEngine_Set-11           	  501634	      2057 ns/op	     533 B/op	      33 allocs/op
// BenchmarkShardedEngine_Set-11    	 2121034	       558.9 ns/op	     529 B/op	      33 allocs/op
// BenchmarkEngine_Del-11           	24150943	        41.55 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	13812578	        72.77 ns/op	       0 B/op	       0 allocs/op

const (
	alphabet        = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	logShardsAmount = 5
	strsAmount      = 10000
	strLen          = 50
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
				_, _ = engine.Get(strs[i%strsAmount])
			}
		}()
	}

	wg.Wait()
}

func BenchmarkShardedEngine_Get(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	strs := genStrings(strsAmount, strLen)
	engine := shardedengine.NewShardedEngine(logShardsAmount, func() shardedengine.BaseEngine {
		return baseengine.NewEngine()
	})

	fillEngine(engine, strs)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				_, _ = engine.Get(strs[i%strsAmount])
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
	engine := shardedengine.NewShardedEngine(logShardsAmount, func() shardedengine.BaseEngine {
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
	engine := shardedengine.NewShardedEngine(logShardsAmount, func() shardedengine.BaseEngine {
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
