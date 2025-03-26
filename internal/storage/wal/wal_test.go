package wal

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func compareCommands(command *parser.Command, commands []Command) bool {
	for i := range commands {
		if !reflect.DeepEqual(command.Args, commands[i].Args) || int(command.Type) != commands[i].CommandType {
			return false
		}
	}
	return true
}

func newTestWal(t *testing.T, ctx context.Context, batchSize int, batchTimeout time.Duration) (*Wal, *MockLogsReader, *MockLogsWriter) {
	logsReader := NewMockLogsReader(t)
	logsWriter := NewMockLogsWriter(t)

	cfg := &config.Config{
		Wal: &config.Wal{
			FlushingBatchSize:    batchSize,
			FlushingBatchTimeout: batchTimeout,
		},
	}

	wal, err := NewWal(cfg, logsReader, logsWriter, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	go wal.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	return wal, logsReader, logsWriter
}

func TestSave_Timeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal, _, logsWriter := newTestWal(t, ctx, 10, 500*time.Millisecond)
	const commandsNumber = 5
	command := &parser.Command{
		Type: 1,
		Args: []string{"name", "Daniil"},
	}

	var commandsCount atomic.Int32
	logsWriter.EXPECT().Write(mock.MatchedBy(func(commands []Command) bool {
		commandsCount.Add(int32(len(commands)))
		return compareCommands(command, commands)
	})).Return(nil)

	wg := sync.WaitGroup{}
	wg.Add(commandsNumber)

	for range commandsNumber {
		go func() {
			defer wg.Done()
			res := wal.Save(command)
			assert.True(t, res)
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(commandsNumber), commandsCount.Load())
}

func TestSave_BatchOverflow(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const batchSize = 10
	wal, _, logsWriter := newTestWal(t, ctx, batchSize, time.Hour)
	command := &parser.Command{
		Type: 1,
		Args: []string{"name", "Daniil"},
	}

	logsWriter.EXPECT().Write(mock.MatchedBy(func(commands []Command) bool {
		return len(commands) == batchSize && compareCommands(command, commands)
	})).Return(nil).Once()

	wg := sync.WaitGroup{}
	wg.Add(batchSize)

	for range batchSize {
		go func() {
			defer wg.Done()
			res := wal.Save(command)
			assert.True(t, res)
		}()
	}

	wg.Wait()
}

func TestSave_Error(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal, _, logsWriter := newTestWal(t, ctx, 10, 500*time.Millisecond)
	logsWriter.EXPECT().Write(mock.Anything).Return(errors.New("write error")).Once()

	res := wal.Save(&parser.Command{
		Type: 1,
		Args: []string{"name", "Daniil"},
	})
	assert.False(t, res)
}

func TestSave_ContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal, _, logsWriter := newTestWal(t, ctx, 10, time.Hour)
	command := &parser.Command{
		Type: 1,
		Args: []string{"name", "Daniil"},
	}

	logsWriter.EXPECT().Write(mock.MatchedBy(func(commands []Command) bool {
		return len(commands) == 1 && compareCommands(command, commands)
	})).Return(nil).Once()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		res := wal.Save(command)
		assert.True(t, res)
	}()
	time.Sleep(100 * time.Millisecond)

	cancel()
	wg.Wait()
}

func TestRecover_Success(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wal, logsReader, _ := newTestWal(t, ctx, 10, 500*time.Millisecond)
	commands := []Command{
		{LSN: 3, CommandType: 0, Args: []string{"name"}},
		{LSN: 2, CommandType: 1, Args: []string{"name", "Daniil"}},
		{LSN: 1, CommandType: 2, Args: []string{"name"}},
	}
	logsReader.EXPECT().Read().Return(commands, nil).Once()
	slices.SortFunc(commands, func(command1 Command, command2 Command) int {
		return command1.LSN - command2.LSN
	})

	res, err := wal.Recover()
	require.NoError(t, err)

	assert.Len(t, res, 3)
	for i := range res {
		assert.Equal(t, commands[i].Args, res[i].Args)
		assert.Equal(t, commands[i].CommandType, int(res[i].Type))
	}
	assert.Equal(t, wal.batch.lsn, 4)
}

func TestRecover_Error(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wal, logsReader, _ := newTestWal(t, ctx, 10, 500*time.Millisecond)
	logsReader.EXPECT().Read().Return(nil, errors.New("recover error")).Once()

	_, err := wal.Recover()
	assert.Error(t, err)
}
