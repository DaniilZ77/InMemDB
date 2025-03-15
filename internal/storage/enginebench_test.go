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
// BenchmarkEngine_Get-11           	 1085605	      1100 ns/op	     178 B/op	      11 allocs/op
// BenchmarkShardedEngine_Get-11    	 1885983	       632.5 ns/op	     177 B/op	      11 allocs/op
// BenchmarkEngine_Set-11           	  628112	      1756 ns/op	       3 B/op	       0 allocs/op
// BenchmarkShardedEngine_Set-11    	 1632589	       745.4 ns/op	       1 B/op	       0 allocs/op
// BenchmarkEngine_Del-11           	 1379066	       858.4 ns/op	       1 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	 2951918	       410.5 ns/op	       0 B/op	       0 allocs/op

// logShardsAmount = 4
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	 1073739	      1089 ns/op	     178 B/op	      11 allocs/op
// BenchmarkShardedEngine_Get-11    	 2763898	       431.3 ns/op	     176 B/op	      11 allocs/op
// BenchmarkEngine_Set-11           	  643636	      1631 ns/op	       3 B/op	       0 allocs/op
// BenchmarkShardedEngine_Set-11    	 2295094	       524.7 ns/op	       1 B/op	       0 allocs/op
// BenchmarkEngine_Del-11           	 1389600	       852.1 ns/op	       1 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	 4220278	       300.9 ns/op	       0 B/op	       0 allocs/op

// logShardsAmount = 5
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	 1033834	      1164 ns/op	     178 B/op	      11 allocs/op
// BenchmarkShardedEngine_Get-11    	 3601521	       327.4 ns/op	     176 B/op	      11 allocs/op
// BenchmarkEngine_Set-11           	  622088	      1633 ns/op	       3 B/op	       0 allocs/op
// BenchmarkShardedEngine_Set-11    	 2852871	       403.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkEngine_Del-11           	 1394392	       853.2 ns/op	       1 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	 5743671	       206.1 ns/op	       0 B/op	       0 allocs/op

// logShardsAmount = 6
// strsAmount      = 10000
// strLen          = 50
// BenchmarkEngine_Get-11           	 1046714	      1120 ns/op	     178 B/op	      11 allocs/op
// BenchmarkShardedEngine_Get-11    	 3965648	       286.7 ns/op	     176 B/op	      11 allocs/op
// BenchmarkEngine_Set-11           	  630590	      1629 ns/op	       3 B/op	       0 allocs/op
// BenchmarkShardedEngine_Set-11    	 3589161	       342.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkEngine_Del-11           	 1400953	       864.4 ns/op	       1 B/op	       0 allocs/op
// BenchmarkShardedEngine_Del-11    	 6024118	       175.2 ns/op	       0 B/op	       0 allocs/op

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
	engine := shardedengine.NewShardedEngine(logShardsAmount, func() shardedengine.BaseEngine {
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
