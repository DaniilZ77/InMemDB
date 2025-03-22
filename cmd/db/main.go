package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/DaniilZ77/InMemDB/internal/app"
	"github.com/DaniilZ77/InMemDB/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewConfig()
	log := newLogger(cfg)

	app := app.NewApp(ctx, cfg, log)
	go func() {
		if err := app.Run(); err != nil {
			panic("failed to run app: " + err.Error())
		}
	}()

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	<-interruptCh
	cancel()
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
