package engine

import (
	"math/rand/v2"
	"runtime"
	"sync"
	"testing"
)

func genString() string {
	const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var res string
	for range rand.IntN(50) + 1 {
		res += string(alpha[rand.IntN(len(alpha))])
	}

	return res
}

type engineSetter interface {
	Set(key, value string)
}

func fillEngine(engine engineSetter) {
	for range 10000 {
		engine.Set(genString(), genString())
	}
}

func BenchmarkEngine_Get(b *testing.B) {
	engine, _ := NewEngine(16)
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	fillEngine(engine)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Get(genString())
			}
		}()
	}

	wg.Wait()
}

func BenchmarkEngine_Set(b *testing.B) {
	engine, _ := NewEngine(16)
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	fillEngine(engine)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Set(genString(), genString())
			}
		}()
	}

	wg.Wait()
}

func BenchmarkEngine_Del(b *testing.B) {
	engine, _ := NewEngine(16)
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	fillEngine(engine)

	for range runtime.NumCPU() {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				engine.Del(genString())
			}
		}()
	}

	wg.Wait()
}
