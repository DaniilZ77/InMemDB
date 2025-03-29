package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage/replication"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Ctx           context.Context
	mainServer    *server.Server
	replicaServer *server.Server
}

func NewApp(ctx context.Context, config *config.Config) (*App, error) {
	group, groupCtx := errgroup.WithContext(ctx)

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

	wal, replica, err := NewWalReplica(groupCtx, config, log)
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

	group.Go(func() error {
		return mainServer.Run(groupCtx, func(b []byte) ([]byte, error) {
			response := database.Execute(string(b))
			return []byte(response), nil
		})
	})

	var replicaServer *server.Server
	switch r := replica.(type) {
	case *replication.Master:
		replicaServer, err = NewReplicaServer(config, log)
		if err != nil {
			return nil, fmt.Errorf("failed to init replica server: %w", err)
		}

		group.Go(func() error {
			return replicaServer.Run(groupCtx, func(b []byte) ([]byte, error) {
				response, err := r.HandleRequest(b)
				return response, err
			})
		})
	case *replication.Slave:
		group.Go(func() error {
			err = r.Start(groupCtx)
			log.Info("slave replica stopped", slog.Any("error", err))
			return err
		})
	}

	return &App{
		Ctx:           groupCtx,
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
