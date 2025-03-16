package storage

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
)

type Compute interface {
	Parse(source string) (*parser.Command, error)
}

type Engine interface {
	Del(key string)
	Get(key string) (string, error)
	Set(key, value string)
}

type Database struct {
	compute Compute
	engine  Engine
	log     *slog.Logger
}

func NewDatabase(compute Compute, engine Engine, log *slog.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("compute is nil")
	}
	if engine == nil {
		return nil, errors.New("engine is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	return &Database{
		compute: compute,
		engine:  engine,
		log:     log,
	}, nil
}

func (d *Database) Execute(source string) string {
	command, err := d.compute.Parse(source)
	if err != nil {
		return fmt.Sprintf("ERROR(%s)", err.Error())
	}

	switch command.Type {
	case parser.SET:
		d.engine.Set(command.Args[0], command.Args[1])
		return "OK"
	case parser.GET:
		res, err := d.engine.Get(command.Args[0])
		if err != nil {
			if errors.Is(err, engine.ErrKeyNotFound) {
				return "NIL"
			}
			return "ERROR(internal error)"
		}
		return res
	case parser.DEL:
		d.engine.Del(command.Args[0])
		return "OK"
	}

	return "ERROR(internal error)"
}
