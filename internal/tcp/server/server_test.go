package server

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/tcp/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	database := mocks.NewDatabase(t)

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

	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
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
		defer conn.Close()

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
			defer conn.Close()
		}

		database.EXPECT().Execute(mock.MatchedBy(func(source string) bool {
			return strings.TrimSpace(source) == command
		})).Return("OK").Once()

		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)

		err = conn.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, err)
		defer conn.Close()

		_, err = fmt.Fprintln(conn, command)
		require.NoError(t, err)

		_, err = io.ReadAll(conn)
		assert.Error(t, err)
	})

	t.Run("exceed read deadline", func(t *testing.T) {
		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)
		defer conn.Close()

		time.Sleep(5 * time.Second)

		body, _ := io.ReadAll(conn)
		assert.Empty(t, body)
	})
}
