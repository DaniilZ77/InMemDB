package storage

import (
	"errors"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
)

type Parser interface {
	Parse(source string) (*parser.Command, error)
}

type Engine interface {
	Del(key string)
	Get(key string) (*string, error)
	Set(key, value string)
}

type Database struct {
	parser Parser
	engine Engine
	log    *slog.Logger
}

func NewDatabase(parser Parser, engine Engine, log *slog.Logger) *Database {
	if parser == nil {
		panic("parser is nil")
	}
	if engine == nil {
		panic("engine is nil")
	}
	if log == nil {
		panic("logger is nil")
	}

	return &Database{
		parser: parser,
		engine: engine,
		log:    log,
	}
}

func (d *Database) Execute(source string) string {
	command, err := d.parser.Parse(source)
	if err != nil {
		return "ERROR"
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
			return "ERROR"
		}
		return *res
	case parser.DEL:
		d.engine.Del(command.Args[0])
		return "OK"
	}

	return "ERROR"
}
