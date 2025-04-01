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

	database, err := NewDatabase(compute, engine, wal, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
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
				engine.EXPECT().Get("name").Return("Daniil", true).Once()
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
				engine.EXPECT().Set("name", "Daniil").Return().Once()
				wal.EXPECT().Save(command).Return(true).Once()
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
				engine.EXPECT().Del("name").Return().Once()
				wal.EXPECT().Save(command).Return(true).Once()
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
				engine.EXPECT().Get("name").Return("", false).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res := database.Execute(tt.command)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestExecute_ParserError(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)

	database, err := NewDatabase(compute, engine, nil, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	compute.EXPECT().Parse("get name").Return(nil, errors.New("internal error")).Once()

	res := database.Execute("get name")
	assert.Contains(t, res, "ERROR")
}

func TestExecute_NilWal(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)

	database, err := NewDatabase(compute, engine, nil, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	commandStr := "set name Daniil"
	compute.EXPECT().Parse(commandStr).Return(&parser.Command{
		Type: parser.SET,
		Args: []string{"name", "Daniil"},
	}, nil)
	engine.EXPECT().Set("name", "Daniil").Return().Once()

	res := database.Execute(commandStr)
	assert.Equal(t, "OK", res)
}

func TestExecute_WalSaveError(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	wal := NewMockWal(t)

	database, err := NewDatabase(compute, engine, wal, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	command := &parser.Command{
		Type: parser.SET,
		Args: []string{"name", "Daniil"},
	}
	commandStr := "set name Daniil"
	compute.EXPECT().Parse(commandStr).Return(command, nil).Once()
	wal.EXPECT().Save(command).Return(false).Once()

	res := database.Execute(commandStr)
	assert.Contains(t, res, "ERROR")
}

func TestRecover_Success(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	w := NewMockWal(t)

	database, err := NewDatabase(compute, engine, w, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	w.EXPECT().Recover().Return([]wal.Command{
		{CommandType: 1, Args: []string{"name", "Daniil"}},
		{CommandType: 2, Args: []string{"name"}},
		{CommandType: 0, Args: []string{"name"}},
	}, nil).Once()
	engine.EXPECT().Set("name", "Daniil").Return().Once()
	engine.EXPECT().Del("name").Return().Once()

	err = database.Recover()
	assert.Nil(t, err)
}

func TestRecover_NilWal(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)

	database, err := NewDatabase(compute, engine, nil, nil, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	err = database.Recover()
	assert.Nil(t, err)
}

func TestSlaveReplica_ForbiddenCommands(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	replica := NewMockReplication(t)

	replica.EXPECT().IsSlave().Return(true)
	replica.EXPECT().GetReplicationStream().Return(nil).Once()

	database, err := NewDatabase(compute, engine, nil, replica, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	compute.EXPECT().Parse(mock.Anything).Return(&parser.Command{Type: parser.SET}, nil).Once()
	res := database.Execute("set a b")
	assert.Equal(t, errReplicaNotSupport, res)

	compute.EXPECT().Parse(mock.Anything).Return(&parser.Command{Type: parser.DEL}, nil).Once()
	res = database.Execute("del a")
	assert.Equal(t, errReplicaNotSupport, res)

	time.Sleep(100 * time.Millisecond)
}

func TestReplicationStream(t *testing.T) {
	t.Parallel()

	compute := NewMockCompute(t)
	engine := NewMockEngine(t)
	replica := NewMockReplication(t)

	engine.EXPECT().Set(mock.Anything, mock.Anything).Return().Once()
	engine.EXPECT().Del(mock.Anything).Return().Once()

	replicationStream := make(chan []wal.Command)
	replica.EXPECT().IsSlave().Return(true).Once()
	replica.EXPECT().GetReplicationStream().Return(replicationStream).Once()

	_, err := NewDatabase(compute, engine, nil, replica, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	commands := []wal.Command{
		{CommandType: setCommand, Args: []string{"name", "Daniil"}},
		{CommandType: delCommand, Args: []string{"name"}},
	}

	go func() {
		replicationStream <- commands
	}()

	time.Sleep(100 * time.Millisecond)
}
