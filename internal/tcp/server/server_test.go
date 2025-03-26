package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	database := NewMockDatabase(t)

	cfg := &config.Config{
		Network: config.Network{
			Address:        "127.0.0.1:0",
			MaxConnections: 5,
			MaxMessageSize: 100,
			IdleTimeout:    5 * time.Second,
		},
	}

	server, err := NewServer(cfg, database, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		server.Run(ctx) // nolint
	}()

	time.Sleep(500 * time.Millisecond)

	address := server.lst.Addr().String()
	command := "set name Daniil"

	t.Run("success", func(t *testing.T) {
		database.EXPECT().Execute(mock.MatchedBy(func(source string) bool {
			return strings.TrimSpace(source) == command
		})).Return("OK").Once()

		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		defer conn.Close() // nolint

		resp := bufio.NewReader(conn)

		_, err = fmt.Fprintln(conn, command)
		require.NoError(t, err)

		body, err := resp.ReadString('\n')
		require.NoError(t, err)

		assert.Equal(t, "OK", strings.TrimSpace(body))
	})

	t.Run("exceed client limit", func(t *testing.T) {
		for range cfg.Network.MaxConnections {
			conn, err := net.Dial("tcp", address)
			require.NoError(t, err)
			defer conn.Close() // nolint
		}

		database.EXPECT().Execute(mock.MatchedBy(func(source string) bool {
			return strings.TrimSpace(source) == command
		})).Return("OK").Once()

		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)

		err = conn.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, err)
		defer conn.Close() // nolint

		_, err = fmt.Fprintln(conn, command)
		require.NoError(t, err)

		_, err = io.ReadAll(conn)
		assert.Error(t, err)
	})

	t.Run("exceed read deadline", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		defer conn.Close() // nolint

		time.Sleep(5 * time.Second)

		body, _ := io.ReadAll(conn)
		assert.Empty(t, body)
	})

	t.Run("context cancel", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		defer conn.Close() // nolint

		cancel()

		cancelErr := errors.New("context was cancelled")
		ctx, cancel := context.WithCancelCause(context.Background())
		done := make(chan struct{}, 1)
		go func() {
			server.Shutdown(ctx)
			done <- struct{}{}
		}()
		cancel(cancelErr)

		select {
		case <-time.After(time.Second):
		case <-done:
		}

		assert.ErrorIs(t, cancelErr, context.Cause(ctx))
	})
}
