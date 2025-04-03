package app

import (
	"errors"
	"log/slog"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/tcp/client"
)

const defaultRetries = 10

func NewClient(address string, log *slog.Logger, opts ...client.ClientOption) (*client.Client, error) {
	retries := defaultRetries
	for retries > 0 {
		client, err := client.NewClient(address, opts...)
		if err == nil {
			return client, nil
		}

		log.Info("client failed to connect to server",
			slog.String("address", address),
			slog.String("error", err.Error()),
			slog.Int("retries", retries),
		)

		retries--

		time.Sleep(500 * time.Millisecond)
	}

	log.Error("client failed to connect to server", slog.String("address", address))

	return nil, errors.New("client failed to connect to server")
}
