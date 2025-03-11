package app

import (
	"log/slog"

	compute "github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server"
)

type App struct {
	server *server.Server
}

func New(
	cfg *config.Config,
	log *slog.Logger) *App {

	parser := compute.NewParser(log)
	engine := engine.NewEngine()
	database := storage.NewDatabase(parser, engine, log)
	server, err := server.New(cfg, database, log)
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
