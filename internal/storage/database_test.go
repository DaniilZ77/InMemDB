package storage

import (
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecute_Success(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	wal := NewMockWal(t)
	coordinator := NewMockCoordinator(t)

	database, err := NewDatabase(compute, engine, coordinator, wal, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	tests := []struct {
		name     string
		command  string
		expected string
		mock     func()
	}{
		{
			name:     "get command",
			command:  "get name",
			expected: "Daniil",
			mock: func() {
				compute.EXPECT().Parse("get name").Return(&parser.Command{
					Type: parser.GET,
					Args: []string{"name"},
				}, nil).Once()
				coordinator.EXPECT().Get("name").Return("Daniil", true).Once()
			},
		},
		{
			name:     "set command",
			command:  "set name Daniil",
			expected: "OK",
			mock: func() {
				command := &parser.Command{
					Type: parser.SET,
					Args: []string{"name", "Daniil"},
				}
				compute.EXPECT().Parse("set name Daniil").Return(command, nil).Once()
				coordinator.EXPECT().Set("name", "Daniil").Return(nil).Once()
			},
		},
		{
			name:     "del command",
			command:  "del name",
			expected: "OK",
			mock: func() {
				command := &parser.Command{
					Type: parser.DEL,
					Args: []string{"name"},
				}
				compute.EXPECT().Parse("del name").Return(command, nil).Once()
				coordinator.EXPECT().Del("name").Return(nil).Once()
			},
		},
		{
			name:     "key not found",
			command:  "get name",
			expected: "NIL",
			mock: func() {
				compute.EXPECT().Parse("get name").Return(&parser.Command{
					Type: parser.GET,
					Args: []string{"name"},
				}, nil).Once()
				coordinator.EXPECT().Get("name").Return("", false).Once()
			},
		},
	}

	const (
		client = "127.0.0.1:5432"
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res := database.Execute(client, tt.command)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestExecute_ParserError(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	coordinator := NewMockCoordinator(t)

	database, err := NewDatabase(compute, engine, coordinator, nil, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	compute.EXPECT().Parse("get name").Return(nil, errors.New("internal error")).Once()

	const (
		client = "127.0.0.1:5432"
	)
	res := database.Execute(client, "get name")
	assert.Contains(t, res, "ERROR")
}

func TestRecover_Success(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	w := NewMockWal(t)
	coordinator := NewMockCoordinator(t)

	database, err := NewDatabase(compute, engine, coordinator, w, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	var txID int64 = 1
	w.EXPECT().Recover().Return([]wal.Command{
		{CommandType: setCommand, TxID: txID, Args: []string{"name", "Daniil"}},
		{CommandType: delCommand, TxID: txID, Args: []string{"name"}},
		{CommandType: 0, TxID: txID, Args: []string{"name"}},
		{CommandType: commitCommand, TxID: txID},
	}, nil).Once()
	engine.EXPECT().Set(mock.Anything, "name", mock.Anything).Return().Once()
	engine.EXPECT().Set(mock.Anything, "name", mock.Anything).Return().Once()

	err = database.Recover()
	assert.Nil(t, err)
}

func TestRecover_NilWal(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	coordinator := NewMockCoordinator(t)

	database, err := NewDatabase(compute, engine, coordinator, nil, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	err = database.Recover()
	assert.Nil(t, err)
}

func TestSlaveReplica_ForbiddenCommands(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	replica := NewMockReplication(t)
	coordinator := NewMockCoordinator(t)

	replica.EXPECT().IsSlave().Return(true)
	replica.EXPECT().GetReplicationStream().Return(nil).Once()

	database, err := NewDatabase(compute, engine, coordinator, nil, replica, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	const (
		client = "127.0.0.1:5432"
	)
	compute.EXPECT().Parse(mock.Anything).Return(&parser.Command{Type: parser.SET}, nil).Once()
	res := database.Execute(client, "set a b")
	assert.Equal(t, errReplicaNotSupport, res)

	compute.EXPECT().Parse(mock.Anything).Return(&parser.Command{Type: parser.DEL}, nil).Once()
	res = database.Execute(client, "del a")
	assert.Equal(t, errReplicaNotSupport, res)

	time.Sleep(100 * time.Millisecond)
}

func TestReplicationStream(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	replica := NewMockReplication(t)
	coordinator := NewMockCoordinator(t)

	engine.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything).Return().Once()
	engine.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything).Return().Once()

	replicationStream := make(chan []wal.Command)
	replica.EXPECT().IsSlave().Return(true).Once()
	replica.EXPECT().GetReplicationStream().Return(replicationStream).Once()

	_, err := NewDatabase(compute, engine, coordinator, nil, replica, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	var txID int64 = 1
	commands := []wal.Command{
		{CommandType: setCommand, TxID: txID, Args: []string{"name", "Daniil"}},
		{CommandType: delCommand, TxID: txID, Args: []string{"name"}},
		{CommandType: commitCommand, TxID: txID},
	}

	go func() {
		replicationStream <- commands
	}()

	time.Sleep(100 * time.Millisecond)
}
