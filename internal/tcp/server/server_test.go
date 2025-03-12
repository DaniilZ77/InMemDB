package server

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/storage"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAddress = "127.0.0.1:3223"
)

func TestMain(m *testing.M) {
	log := slog.New(slog.NewJSONHandler(io.Discard, nil))

	cfg := &config.Config{
		Network: config.Network{
			Address:        testAddress,
			MaxConnections: 3,
			MaxMessageSize: 200,
			IdleTimeout:    5 * time.Minute,
		},
	}

	compute, err := parser.NewParser(log)
	if err != nil {
		panic(err)
	}

	engine := engine.NewEngine()

	database, err := storage.NewDatabase(compute, engine, log)
	if err != nil {
		panic(err)
	}

	server, err := NewServer(cfg, database, log)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	os.Exit(m.Run())
}

func TestSet_Success(t *testing.T) {
	conn, err := net.Dial("tcp", testAddress)
	require.NoError(t, err)

	defer conn.Close()

	server := bufio.NewReader(conn)

	_, err = fmt.Fprintln(conn, "set name Daniil")
	require.NoError(t, err)

	resp, err := server.ReadString('\n')
	require.NoError(t, err)

	assert.Equal(t, "OK\n", string(resp))
}

func Test_Fail(t *testing.T) {
	conn, err := net.Dial("tcp", testAddress)
	require.NoError(t, err)

	defer conn.Close()

	server := bufio.NewReader(conn)

	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "bad amount of args",
			command: "get a b c",
		},
		{
			name:    "bad amount of args",
			command: "set",
		},
		{
			name:    "bad amount of args",
			command: "del a b",
		},
		{
			name:    "bad command type",
			command: "st name Daniil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = fmt.Fprintln(conn, tt.command)
			require.NoError(t, err)

			resp, err := server.ReadString('\n')
			require.NoError(t, err)

			assert.Contains(t, string(resp), "ERROR")
		})
	}
}

func TestGet_Success(t *testing.T) {
	conn, err := net.Dial("tcp", testAddress)
	require.NoError(t, err)

	defer conn.Close()

	server := bufio.NewReader(conn)

	_, err = fmt.Fprintln(conn, "set name Daniil")
	require.NoError(t, err)

	resp, err := server.ReadString('\n')
	require.NoError(t, err)

	require.Equal(t, "OK\n", string(resp))

	_, err = fmt.Fprintln(conn, "get name")
	require.NoError(t, err)

	resp, err = server.ReadString('\n')
	require.NoError(t, err)

	assert.Equal(t, "Daniil\n", string(resp))
}

func TestDel_Success(t *testing.T) {
	conn, err := net.Dial("tcp", testAddress)
	require.NoError(t, err)

	defer conn.Close()

	server := bufio.NewReader(conn)

	_, err = fmt.Fprintln(conn, "set name Daniil")
	require.NoError(t, err)

	resp, err := server.ReadString('\n')
	require.NoError(t, err)

	require.Equal(t, "OK\n", string(resp))

	_, err = fmt.Fprintln(conn, "del name")
	require.NoError(t, err)

	resp, err = server.ReadString('\n')
	require.NoError(t, err)

	assert.Equal(t, "OK\n", string(resp))

	_, err = fmt.Fprintln(conn, "get name")
	require.NoError(t, err)

	resp, err = server.ReadString('\n')
	require.NoError(t, err)

	assert.Equal(t, "NIL\n", string(resp))
}
