package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage/replication"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server"
)

type App struct {
	mainServer    *server.Server
	replicaServer *server.Server
}

func NewApp(ctx context.Context, config *config.Config) (*App, error) {
	log, err := NewLogger(config)
	if err != nil {
		return nil, err
	}

	parser, err := parser.NewParser(log)
	if err != nil {
		return nil, fmt.Errorf("failed to init parser: %w", err)
	}

	engine, err := NewEngine(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init engine: %w", err)
	}

	wal, replica, err := NewWalReplica(ctx, config, log)
	if err != nil {
		return nil, fmt.Errorf("failed to init wal and replica: %w", err)
	}

	database, err := NewDatabase(parser, engine, wal, replica, log)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	err = database.Recover()
	if err != nil {
		return nil, fmt.Errorf("failed to recover database: %w", err)
	}

	mainServer, err := NewServer(config, log)
	if err != nil {
		return nil, fmt.Errorf("failed to init main server: %w", err)
	}

	go func() {
		if err := mainServer.Run(ctx, func(b []byte) ([]byte, error) {
			response := database.Execute(string(b))
			return []byte(response), nil
		}); err != nil {
			log.Warn("main server stopped", slog.Any("error", err))
		}
	}()

	var replicaServer *server.Server
	switch r := replica.(type) {
	case *replication.Master:
		replicaServer, err = NewReplicaServer(config, log)
		if err != nil {
			return nil, fmt.Errorf("failed to init replica server: %w", err)
		}

		go func() {
			if err := replicaServer.Run(ctx, func(b []byte) ([]byte, error) {
				response, err := r.HandleRequest(b)
				return response, err
			}); err != nil {
				log.Warn("replica server stopped", slog.Any("error", err))
			}
		}()
	case *replication.Slave:
		go func() {
			if err := r.Start(ctx); err != nil {
				log.Warn("slave replica stopped", slog.Any("error", err))
				panic(err)
			}
		}()
	default:
	}

	return &App{
		mainServer:    mainServer,
		replicaServer: replicaServer,
	}, nil
}

func (a *App) Shutdown(ctx context.Context) {
	a.mainServer.Shutdown(ctx)
	if a.replicaServer != nil {
		a.replicaServer.Shutdown(ctx)
	}
}
