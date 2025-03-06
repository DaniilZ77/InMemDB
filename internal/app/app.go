package app

import (
	compute "github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp"
)

type App struct {
	server *tcp.Server
}

func New(
	cfg *config.Config,
	log *slog.Logger) *App {

	parser := compute.NewParser()
	server, err := tcp.New(cfg, log, parser, nil)
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
