package app

import (
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"

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

	engine := engine.NewEngine()
	database, err := storage.NewDatabase(parser, engine, log)
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
