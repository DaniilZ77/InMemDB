package storage

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
)

const (
	setCommand           = 1
	delCommand           = 2
	errReplicaNotSupport = "ERROR(invalid command: replica support only get commands)"
	errInternal          = "ERROR(internal error)"
)

//go:generate mockery --name=Compute --case=snake --inpackage --inpackage-suffix --with-expecter
type Compute interface {
	Parse(source string) (*parser.Command, error)
}

//go:generate mockery --name=Engine --case=snake --inpackage --inpackage-suffix --with-expecter
type Engine interface {
	Del(key string)
	Get(key string) (string, bool)
	Set(key, value string)
}

//go:generate mockery --name=Wal --case=snake --inpackage --inpackage-suffix --with-expecter
type Wal interface {
	Save(command *parser.Command) bool
	Recover() ([]wal.Command, error)
}

//go:generate mockery --name=Replication --case=snake --inpackage --inpackage-suffix --with-expecter
type Replication interface {
	IsSlave() bool
	GetReplicationStream() <-chan []wal.Command
}

type Database struct {
	compute Compute
	engine  Engine
	wal     Wal
	replica Replication
	log     *slog.Logger
}

func NewDatabase(
	compute Compute,
	engine Engine,
	wal Wal,
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

	database := &Database{
		compute: compute,
		engine:  engine,
		wal:     wal,
		replica: replica,
		log:     log,
	}

	if replica != nil && replica.IsSlave() {
		go func() {
			for commands := range replica.GetReplicationStream() {
				database.executeWalCommands(commands)
			}
		}()
	}

	return database, nil
}

func (d *Database) Execute(source string) string {
	command, err := d.compute.Parse(source)
	if err != nil {
		return fmt.Sprintf("ERROR(%s)", err.Error())
	}

	switch command.Type {
	case parser.SET:
		return d.setCommand(command)
	case parser.GET:
		return d.getCommand(command)
	case parser.DEL:
		return d.delCommand(command)
	}

	return errInternal
}

func (d *Database) executeWalCommands(commands []wal.Command) {
	for _, command := range commands {
		switch command.CommandType {
		case setCommand:
			d.engine.Set(command.Args[0], command.Args[1])
		case delCommand:
			d.engine.Del(command.Args[0])
		default:
			d.log.Warn("command type must be one of set or del")
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
	d.executeWalCommands(commands)
	return nil
}

func (d *Database) setCommand(command *parser.Command) string {
	if d.replica != nil && d.replica.IsSlave() {
		return errReplicaNotSupport
	}

	if d.wal == nil || d.wal.Save(command) {
		d.engine.Set(command.Args[0], command.Args[1])
		return "OK"
	}

	return errInternal
}

func (d *Database) getCommand(command *parser.Command) string {
	res, ok := d.engine.Get(command.Args[0])
	if !ok {
		return "NIL"
	}

	return res
}

func (d *Database) delCommand(command *parser.Command) string {
	if d.replica != nil && d.replica.IsSlave() {
		return errReplicaNotSupport
	}

	if d.wal == nil || d.wal.Save(command) {
		d.engine.Del(command.Args[0])
		return "OK"
	}

	return errInternal
}
