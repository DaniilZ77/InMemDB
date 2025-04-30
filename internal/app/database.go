package app

import (
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
	"github.com/DaniilZ77/InMemDB/internal/storage/mvcc"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
)

func NewDatabase(
	parser *parser.Parser,
	engine *engine.Engine,
	coordinator *mvcc.Coordinator,
	wal *wal.Wal,
	replica any,
	log *slog.Logger,
) (database *storage.Database, err error) {
	if wal == nil {
		return storage.NewDatabase(parser, engine, coordinator, nil, nil, log)
	}

	if replica == nil {
		return storage.NewDatabase(parser, engine, coordinator, wal, nil, log)
	}

	return storage.NewDatabase(parser, engine, coordinator, wal, replica.(storage.Replication), log)
}
