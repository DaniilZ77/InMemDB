package app

import (
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
	"github.com/DaniilZ77/InMemDB/internal/storage/sharded"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server"
)

type App struct {
	server *server.Server
}

func NewApp(
	cfg *config.Config,
	log *slog.Logger) *App {

	parser, err := parser.NewParser(log)
	if err != nil {
		panic("failed to init parser: " + err.Error())
	}

	var eng storage.Engine
	if cfg.Engine.LogShardsAmount == 0 {
		eng = engine.NewEngine()
	} else if cfg.Engine.LogShardsAmount > 0 {
		eng = sharded.NewShardedEngine(cfg.Engine.LogShardsAmount, func() sharded.BaseEngine {
			return engine.NewEngine()
		})
	} else {
		panic("shards amount must be greater than 0")
	}

	database, err := storage.NewDatabase(parser, eng, log)
	if err != nil {
		panic("failed to init database: " + err.Error())
	}

	server, err := server.NewServer(cfg, database, log)
	if err != nil {
		panic("failed to init server: " + err.Error())
	}

	return &App{server}
}

func (a *App) Run() error {
	if err := a.server.Run(); err != nil {
		return err
	}

	return nil
}
