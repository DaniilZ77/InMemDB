package concurrency

import "sync"

func WithLock(mu sync.Locker, action func()) {
	mu.Lock()
	defer mu.Unlock()
	action()
}
