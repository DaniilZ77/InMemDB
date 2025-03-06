package main

import (
	"log/slog"
	"os"

	"github.com/DaniilZ77/InMemDB/internal/app"
	"github.com/DaniilZ77/InMemDB/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.New()
	log := newLogger(cfg.Env)

	app := app.New(cfg, log)
	if err := app.Run(); err != nil {
		log.Error("failed to run app", slog.Any("error", err))
	}
}

func newLogger(env string) *slog.Logger {
	var log *slog.Logger

	opts := &slog.HandlerOptions{AddSource: true}

	switch env {
	case envLocal:
		opts.Level = slog.LevelDebug
		log = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case envProd:
		opts.Level = slog.LevelInfo
		log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		panic("unknown env")
	}

	return log
}
