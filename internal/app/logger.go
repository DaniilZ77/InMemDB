package app

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/DaniilZ77/InMemDB/internal/config"
)

func NewLogger(config *config.Config) (*slog.Logger, error) {
	var log *slog.Logger
	opts := &slog.HandlerOptions{AddSource: true}

	switch strings.ToUpper(config.LogLevel) {
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
		return nil, errors.New("inavalid log level")
	}

	return log, nil
}
