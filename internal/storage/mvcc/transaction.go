package mvcc

import (
	"errors"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
)

type Transaction struct {
	beginID     int64
	modified    map[string]*string
	cache       map[string]string
	coordinator *Coordinator
	wal         Wal
	finished    bool
}

var (
	ErrTransactionInterfered = errors.New("transaction interfered")
	ErrTransactionFinished   = errors.New("transaction finished")
	ErrWalFailure            = errors.New("wal failure")
)

func newTransaction(txID int64, coordinator *Coordinator, wal Wal) *Transaction {
	return &Transaction{
		beginID:     txID,
		modified:    make(map[string]*string),
		cache:       make(map[string]string),
		coordinator: coordinator,
		wal:         wal,
	}
}

func (tx *Transaction) Get(key string) (string, bool) {
	if tx.finished {
		return "", false
	}

	if value, ok := tx.modified[key]; ok && value != nil {
		return *value, true
	}
	if value, ok := tx.cache[key]; ok {
		return value, true
	}

	value, ok := tx.coordinator.get(tx.beginID, key)
	if ok {
		tx.cache[key] = value
	}
	return value, ok
}

func (tx *Transaction) Set(key, value string) error {
	if tx.finished {
		return ErrTransactionFinished
	}

	if tx.wal != nil && !tx.wal.Save(tx.beginID, parser.NewSetCommand(key, value)) {
		return ErrWalFailure
	}

	tx.modified[key] = &value
	return nil
}

func (tx *Transaction) Del(key string) error {
	if tx.finished {
		return ErrTransactionFinished
	}

	if tx.wal != nil && !tx.wal.Save(tx.beginID, parser.NewDelCommand(key)) {
		return ErrWalFailure
	}

	tx.modified[key] = nil
	return nil
}

func (tx *Transaction) Commit() error {
	if tx.finished {
		return ErrTransactionFinished
	}
	defer func() { tx.finished = true }()

	if !tx.coordinator.apply(tx.beginID, tx.modified) {
		return ErrTransactionInterfered
	}
	return nil
}

func (tx *Transaction) Rollback() error {
	if tx.finished {
		return ErrTransactionFinished
	}
	defer func() { tx.finished = true }()

	if tx.wal != nil && !tx.wal.Save(tx.beginID, parser.NewRollbackCommand()) {
		return ErrWalFailure
	}
	return nil
}
