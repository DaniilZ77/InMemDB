package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
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
			IdleTimeout:    time.Second,
		},
	}
	server, err := NewServer(cfg, database, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go server.Run(ctx) // nolint
	time.Sleep(100 * time.Millisecond)

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
		defer conn.Close() // nolint

		_, err = fmt.Fprintln(conn, command)
		require.NoError(t, err)

		err = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		require.NoError(t, err)

		_, err = conn.Read(make([]byte, 10))
		assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
	})

	t.Run("exceed read deadline", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		defer conn.Close() // nolint

		time.Sleep(1100 * time.Millisecond)

		_, err = conn.Read(make([]byte, 10))
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("force shutdown after context cancel", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		defer conn.Close() // nolint

		cancel()

		shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
		shutdownCancel()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.Shutdown(shutdownCtx)
		}()
		wg.Wait()

		assert.Error(t, ctx.Err())
	})
}
