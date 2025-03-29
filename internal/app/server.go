package app

import (
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server"
)

const (
	defaultAddress        = ":3223"
	defaultMaxMessageSize = 4 << 10
)

func NewServer(config *config.Config, log *slog.Logger) (*server.Server, error) {
	address := defaultAddress
	var maxMessageSize int
	var err error
	opts := []server.ServerOption{}

	if config.Network != nil {
		if config.Network.Address != "" {
			address = config.Network.Address
		}
		if maxMessageSize, err = parseBytes(config.Network.MaxMessageSize); err != nil {
			maxMessageSize = defaultMaxMessageSize
		}
		opts = append(opts, server.WithMaxMessageSize(maxMessageSize))
		if config.Network.IdleTimeout > 0 {
			opts = append(opts, server.WithIdleTimeout(config.Network.IdleTimeout))
		}
		if config.Network.MaxConnections > 0 {
			opts = append(opts, server.WithMaxConnections(config.Network.MaxConnections))
		}
	}

	server, err := server.NewServer(address, maxMessageSize, log, opts...)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func NewReplicaServer(config *config.Config, log *slog.Logger) (*server.Server, error) {
	masterAddress := defaultMasterAddress
	if config.Replication.MasterAddress != "" {
		masterAddress = config.Replication.MasterAddress
	}

	maxMessageSize, err := parseBytes(config.Wal.MaxSegmentSize)
	if err != nil {
		maxMessageSize = defaultMaxMessageSize
	}

	server, err := server.NewServer(masterAddress, maxMessageSize, log)
	if err != nil {
		return nil, err
	}

	return server, nil
}
