package storage

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/concurrency"
	"github.com/DaniilZ77/InMemDB/internal/storage/mvcc"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
)

const (
	setCommand      = 1
	delCommand      = 2
	commitCommand   = 3
	rollbackCommand = 4
	beginCommand    = 5

	errReplicaNotSupport  = "ERROR(invalid command: replica support only get commands)"
	errInternal           = "ERROR(internal error)"
	errInvalidTransaction = "ERROR(invalid transaction)"
)

//go:generate mockery --name=Compute --case=snake --inpackage --inpackage-suffix --with-expecter
type Compute interface {
	Parse(source string) (*parser.Command, error)
}

//go:generate mockery --name=Engine --case=snake --inpackage --inpackage-suffix --with-expecter
type Engine interface {
	Get(version int64, key string) (string, bool)
	Set(version int64, key string, value *string)
}

//go:generate mockery --name=Wal --case=snake --inpackage --inpackage-suffix --with-expecter
type Wal interface {
	Recover() ([]wal.Command, error)
}

type Coordinator interface {
	BeginTransaction() *mvcc.Transaction
	Set(key, value string) error
	Get(key string) (string, bool)
	Del(key string) error
}

type Transaction interface {
	Commit() error
	Rollback() error
	Del(key string) error
	Get(key string) (string, bool)
	Set(key string, value string) error
}

//go:generate mockery --name=Replication --case=snake --inpackage --inpackage-suffix --with-expecter
type Replication interface {
	IsSlave() bool
	GetReplicationStream() <-chan []wal.Command
}

type Database struct {
	compute      Compute
	engine       Engine
	coordinator  Coordinator
	wal          Wal
	replica      Replication
	log          *slog.Logger
	mu           sync.RWMutex
	transactions map[string]Transaction
}

func NewDatabase(
	compute Compute,
	engine Engine,
	coordinator Coordinator,
	w Wal,
	replica Replication,
	log *slog.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("compute is nil")
	}
	if engine == nil {
		return nil, errors.New("engine is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}
	if coordinator == nil {
		return nil, errors.New("coordinator is nil")
	}

	database := &Database{
		compute:      compute,
		engine:       engine,
		coordinator:  coordinator,
		wal:          w,
		replica:      replica,
		log:          log,
		transactions: make(map[string]Transaction),
	}

	if replica != nil && replica.IsSlave() {
		transactions := make(map[int64][]wal.Command)
		go func() {
			for commands := range replica.GetReplicationStream() {
				database.executeWalCommands(transactions, commands)
			}
		}()
	}

	return database, nil
}

func (d *Database) Execute(client string, source string) string {
	command, err := d.compute.Parse(source)
	if err != nil {
		return fmt.Sprintf("ERROR(%s)", err.Error())
	}

	switch command.Type {
	case parser.SET:
		return d.setCommand(client, command)
	case parser.GET:
		return d.getCommand(client, command)
	case parser.DEL:
		return d.delCommand(client, command)
	case parser.BEGIN, parser.COMMIT, parser.ROLLBACK:
		return d.txCommand(client, command)
	}

	return errInternal
}

func (d *Database) executeWalCommands(transactions map[int64][]wal.Command, commands []wal.Command) {
	for _, command := range commands {
		switch command.CommandType {
		case commitCommand:
			for _, v := range transactions[command.TxID] {
				if v.CommandType == setCommand {
					d.engine.Set(0, v.Args[0], &v.Args[1])
				} else if v.CommandType == delCommand {
					d.engine.Set(0, v.Args[0], nil)
				}
			}
			delete(transactions, command.TxID)
		case rollbackCommand:
			delete(transactions, command.TxID)
		case setCommand, delCommand:
			transactions[command.TxID] = append(transactions[command.TxID], command)
		}
	}
}

func (d *Database) Recover() error {
	if d.wal == nil {
		return nil
	}
	commands, err := d.wal.Recover()
	if err != nil {
		return err
	}
	d.executeWalCommands(map[int64][]wal.Command{}, commands)
	return nil
}

func (d *Database) setCommand(client string, command *parser.Command) string {
	if d.replica != nil && d.replica.IsSlave() {
		return errReplicaNotSupport
	}

	var tx Transaction
	var ok bool
	concurrency.WithLock(d.mu.RLocker(), func() {
		tx, ok = d.transactions[client]
	})
	if ok {
		if err := tx.Set(command.Args[0], command.Args[1]); err != nil {
			return errInternal
		}
		return "OK"
	}

	if err := d.coordinator.Set(command.Args[0], command.Args[1]); err != nil {
		if errors.Is(err, mvcc.ErrTransactionInterfered) {
			return fmt.Sprintf("ERROR(%s)", err.Error())
		}
		return errInternal
	}

	return "OK"
}

func (d *Database) getCommand(client string, command *parser.Command) string {
	var tx Transaction
	var ok bool
	concurrency.WithLock(d.mu.RLocker(), func() {
		tx, ok = d.transactions[client]
	})
	if ok {
		if value, ok := tx.Get(command.Args[0]); ok {
			return value
		}
		return "NIL"
	}

	res, ok := d.coordinator.Get(command.Args[0])
	if !ok {
		return "NIL"
	}

	return res
}

func (d *Database) delCommand(client string, command *parser.Command) string {
	if d.replica != nil && d.replica.IsSlave() {
		return errReplicaNotSupport
	}

	var tx Transaction
	var ok bool
	concurrency.WithLock(d.mu.RLocker(), func() {
		tx, ok = d.transactions[client]
	})
	if ok {
		if err := tx.Del(command.Args[0]); err != nil {
			return errInternal
		}
		return "OK"
	}

	if err := d.coordinator.Del(command.Args[0]); err != nil {
		if errors.Is(err, mvcc.ErrTransactionInterfered) {
			return fmt.Sprintf("ERROR(%s)", err.Error())
		}
		return "OK"
	}

	return errInternal
}

func (d *Database) txCommand(client string, command *parser.Command) string {
	switch command.Type {
	case parser.BEGIN:
		tx := d.coordinator.BeginTransaction()
		concurrency.WithLock(&d.mu, func() {
			d.transactions[client] = tx
		})
	case parser.COMMIT, parser.ROLLBACK:
		var tx Transaction
		var ok bool
		defer concurrency.WithLock(&d.mu, func() {
			delete(d.transactions, client)
		})
		concurrency.WithLock(d.mu.RLocker(), func() {
			tx, ok = d.transactions[client]
		})
		if !ok {
			return errInvalidTransaction
		}
		if command.Type == parser.COMMIT {
			if err := tx.Commit(); err != nil {
				return errInternal
			}
		} else {
			if err := tx.Rollback(); err != nil {
				return errInternal
			}
		}
	default:
		return errInternal
	}
	return "OK"
}
