package app

import (
	"context"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/disk"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server"
)

type App struct {
	server *server.Server
}

func NewApp(
	ctx context.Context,
	cfg *config.Config,
	log *slog.Logger) *App {

	parser, err := parser.NewParser(log)
	if err != nil {
		panic("failed to init parser: " + err.Error())
	}

	engine, err := engine.NewEngine(cfg.Engine.ShardsNumber)
	if err != nil {
		panic("failed to init engine: " + err.Error())
	}

	var database *storage.Database
	if cfg.Wal != nil {
		disk := disk.NewDisk(cfg.Wal.DataDirectory, cfg.Wal.MaxSegmentSizeBytes, log)
		logsManager := wal.NewLogsManager(disk, log)

		wal, err := wal.NewWal(cfg.Wal.FlushingBatchSize, cfg.Wal.FlushingBatchTimeout, logsManager, logsManager, log)
		if err != nil {
			panic("failed to init wal: " + err.Error())
		}

		go wal.Start(ctx)

		database, err = storage.NewDatabase(parser, engine, wal, log)
		if err != nil {
			panic("failed to init database: " + err.Error())
		}
	} else {
		database, err = storage.NewDatabase(parser, engine, nil, log)
		if err != nil {
			panic("failed to init database: " + err.Error())
		}
	}

	err = database.Recover()
	if err != nil {
		panic("failed to recover database: " + err.Error())
	}

	server, err := server.NewServer(
		cfg.Network.Address,
		cfg.Network.MaxMessageSizeBytes,
		cfg.Network.IdleTimeout,
		cfg.Network.MaxConnections,
		database,
		log,
	)
	if err != nil {
		panic("failed to init server: " + err.Error())
	}

	return &App{server}
}

func (a *App) Run(ctx context.Context) error {
	if err := a.server.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) {
	a.server.Shutdown(ctx)
}
