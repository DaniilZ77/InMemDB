package storage

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
)

//go:generate mockery --name=Compute --with-expecter
type Compute interface {
	Parse(source string) (*parser.Command, error)
}

//go:generate mockery --name=Engine --with-expecter
type Engine interface {
	Del(key string)
	Get(key string) (string, bool)
	Set(key, value string)
}

type Wal interface {
	Start()
	Save(command *parser.Command) bool
	Recover() ([]parser.Command, error)
}

type Database struct {
	compute Compute
	engine  Engine
	wal     Wal
	log     *slog.Logger
}

func NewDatabase(compute Compute, engine Engine, wal Wal, log *slog.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("compute is nil")
	}
	if engine == nil {
		return nil, errors.New("engine is nil")
	}
	if wal == nil {
		return nil, errors.New("wal is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	database := &Database{
		compute: compute,
		engine:  engine,
		wal:     wal,
		log:     log,
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
		if d.wal.Save(command) {
			d.engine.Set(command.Args[0], command.Args[1])
			return "OK"
		}
	case parser.GET:
		res, ok := d.engine.Get(command.Args[0])
		if !ok {
			return "NIL"
		}
		return res
	case parser.DEL:
		if d.wal.Save(command) {
			d.engine.Del(command.Args[0])
			return "OK"
		}
	}

	return "ERROR(internal error)"
}

func (d *Database) Fill(commands []parser.Command) {
	for _, command := range commands {
		switch command.Type {
		case parser.SET:
			d.engine.Set(command.Args[0], command.Args[1])
		case parser.DEL:
			d.engine.Del(command.Args[0])
		default:
			d.log.Warn("command type must be one of set or del")
		}
	}
}
