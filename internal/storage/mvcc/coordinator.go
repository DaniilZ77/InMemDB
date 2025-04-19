package mvcc

import (
	"errors"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/concurrency"
)

type Engine interface {
	ExistsBetween(beginTxID, endTxID int64, key string) bool
	SetMany(txID int64, modified map[string]*string)
	Set(txID int64, key string, value *string)
	Get(txID int64, key string) (string, bool)
}

type Wal interface {
	Save(txID int64, command *parser.Command) bool
}

type Coordinator struct {
	txID   int64
	engine Engine
	wal    Wal
	log    *slog.Logger
	mu     sync.Mutex
	locks  map[string]*lock
}

type lock struct {
	locked int64
}

func NewCoordinator(engine Engine, wal Wal, log *slog.Logger) (*Coordinator, error) {
	if engine == nil {
		return nil, errors.New("engine is nil")
	}
	return &Coordinator{
		engine: engine,
		wal:    wal,
		log:    log,
		locks:  make(map[string]*lock),
	}, nil
}

func (c *Coordinator) BeginTransaction() *Transaction {
	txID := atomic.AddInt64(&c.txID, 1)
	return newTransaction(txID, c, c.wal)
}

func (c *Coordinator) Set(key, value string) error {
	tx := c.BeginTransaction()
	if err := tx.Set(key, value); err != nil {
		c.log.Error("failed to set value", slog.Any("error", err))
		return err
	}
	if err := tx.Commit(); err != nil {
		c.log.Warn("failed to commit", slog.Any("error", err))
		return err
	}
	return nil
}

func (c *Coordinator) Del(key string) error {
	tx := c.BeginTransaction()
	if err := tx.Del(key); err != nil {
		c.log.Error("failed to delete value", slog.Any("error", err))
		return err
	}
	if err := tx.Commit(); err != nil {
		c.log.Warn("failed to commit", slog.Any("error", err))
		return err
	}
	return nil
}

func (c *Coordinator) Get(key string) (string, bool) {
	txID := atomic.AddInt64(&c.txID, 1)
	return c.engine.Get(txID, key)
}

func (c *Coordinator) get(txID int64, key string) (string, bool) {
	var l *lock
	concurrency.WithLock(&c.mu, func() {
		l = c.locks[key]
	})

	if l != nil {
		ownerTxID := atomic.LoadInt64(&l.locked)
		for ownerTxID != 0 && ownerTxID < txID {
			runtime.Gosched()
			ownerTxID = atomic.LoadInt64(&l.locked)
		}
	}

	return c.engine.Get(txID, key)
}

func (c *Coordinator) apply(beginTxID int64, modified map[string]*string) bool {
	var locks []*lock
	for k := range modified {
		concurrency.WithLock(&c.mu, func() {
			l, ok := c.locks[k]
			if !ok {
				l = &lock{}
				c.locks[k] = l
			}
			locks = append(locks, l)
		})
	}

	if !acquireLocks(beginTxID, locks) {
		return false
	}
	defer releaseLocks(locks)

	endTxID := atomic.AddInt64(&c.txID, 1)
	for k := range modified {
		if c.engine.ExistsBetween(beginTxID, endTxID, k) {
			return false
		}
	}

	if c.wal != nil && !c.wal.Save(beginTxID, parser.NewCommitCommand()) {
		return false
	}

	c.engine.SetMany(endTxID, modified)
	return true
}

func acquireLocks(txID int64, locks []*lock) bool {
	for i, lock := range locks {
		if !atomic.CompareAndSwapInt64(&lock.locked, 0, txID) {
			for j := range i {
				atomic.StoreInt64(&locks[j].locked, 0)
			}
			return false
		}
	}
	return true
}

func releaseLocks(locks []*lock) {
	for _, lock := range locks {
		atomic.StoreInt64(&lock.locked, 0)
	}
}
