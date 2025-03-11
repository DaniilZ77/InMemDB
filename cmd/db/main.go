package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/DaniilZ77/InMemDB/internal/app"
	"github.com/DaniilZ77/InMemDB/internal/config"
)

func main() {
	cfg := config.New()
	log := newLogger(cfg)

	app := app.New(cfg, log)
	if err := app.Run(); err != nil {
		log.Error("failed to run app", slog.Any("error", err))
	}
}

func newLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger
	opts := &slog.HandlerOptions{AddSource: true}

	switch strings.ToUpper(cfg.LogLevel) {
	case slog.LevelDebug.String():
		opts.Level = slog.LevelDebug
		log = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case slog.LevelInfo.String():
		opts.Level = slog.LevelInfo
		log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case slog.LevelWarn.String():
		opts.Level = slog.LevelWarn
		log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case slog.LevelError.String():
		opts.Level = slog.LevelError
		log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		panic("unknown log level")
	}
	return log
}
