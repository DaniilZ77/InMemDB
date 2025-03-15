package app

import (
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/baseengine"
	"github.com/DaniilZ77/InMemDB/internal/storage/shardedengine"

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
		eng = baseengine.NewEngine()
	} else if cfg.Engine.LogShardsAmount > 0 {
		eng = shardedengine.NewShardedEngine(cfg.Engine.LogShardsAmount, func() shardedengine.BaseEngine {
			return baseengine.NewEngine()
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
