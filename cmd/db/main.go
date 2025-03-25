package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/app"
	"github.com/DaniilZ77/InMemDB/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.NewConfig()
	log := newLogger(cfg)

	app := app.NewApp(ctx, cfg, log)
	go func() {
		if err := app.Run(ctx); err != nil {
			log.Error("failed to run app", slog.Any("error", err))
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	app.Shutdown(ctx)
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
