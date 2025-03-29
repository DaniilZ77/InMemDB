package wal

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestWal(t *testing.T, ctx context.Context, batchSize int, batchTimeout time.Duration) (*Wal, *MockLogsReader, *MockLogsWriter) {
	logsReader := NewMockLogsReader(t)
	logsWriter := NewMockLogsWriter(t)

	wal, err := NewWal(batchSize, batchTimeout, logsReader, logsWriter, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	go wal.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	return wal, logsReader, logsWriter
}

func TestSave_Timeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	wal, _, logsWriter := newTestWal(t, ctx, 10, 100*time.Millisecond)
	parserCommands := []parser.Command{
		{Type: parser.SET, Args: []string{"name", "Daniil"}},
		{Type: parser.DEL, Args: []string{"name"}},
		{Type: parser.GET, Args: []string{"name"}},
	}

	var commandsCount atomic.Int32
	logsWriter.EXPECT().WriteLogs(mock.MatchedBy(func(commands []Command) bool {
		commandsCount.Add(int32(len(commands)))
		return true
	})).Return(nil)

	wg := sync.WaitGroup{}
	wg.Add(len(parserCommands))

	for i := range len(parserCommands) {
		go func() {
			defer wg.Done()
			res := wal.Save(&parserCommands[i])
			assert.True(t, res)
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(len(parserCommands)), commandsCount.Load())
}

func TestSave_BatchOverflow(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	const batchSize = 3
	wal, _, logsWriter := newTestWal(t, ctx, batchSize, time.Hour)
	parserCommands := []parser.Command{
		{Type: parser.SET, Args: []string{"name", "Daniil"}},
		{Type: parser.DEL, Args: []string{"name"}},
		{Type: parser.GET, Args: []string{"name"}},
	}

	logsWriter.EXPECT().WriteLogs(mock.MatchedBy(func(commands []Command) bool {
		return len(commands) == batchSize
	})).Return(nil).Once()

	wg := sync.WaitGroup{}
	wg.Add(batchSize)
	for i := range batchSize {
		go func() {
			defer wg.Done()
			res := wal.Save(&parserCommands[i])
			assert.True(t, res)
		}()
	}

	wg.Wait()
}

func TestSave_Error(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	wal, _, logsWriter := newTestWal(t, ctx, 10, 500*time.Millisecond)
	logsWriter.EXPECT().WriteLogs(mock.Anything).Return(errors.New("write error")).Once()

	res := wal.Save(&parser.Command{
		Type: parser.SET,
		Args: []string{"name", "Daniil"},
	})
	assert.False(t, res)
}

func TestSave_ContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	wal, _, logsWriter := newTestWal(t, ctx, 10, time.Hour)
	parserCommands := []parser.Command{{Type: parser.SET, Args: []string{"name", "Daniil"}}}

	logsWriter.EXPECT().WriteLogs(mock.MatchedBy(func(commands []Command) bool {
		return len(commands) == 1
	})).Return(nil).Once()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		res := wal.Save(&parserCommands[0])
		assert.True(t, res)
	}()
	time.Sleep(100 * time.Millisecond)

	cancel()
	wg.Wait()
}

func TestRecover_Success(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	wal, logsReader, _ := newTestWal(t, ctx, 10, 500*time.Millisecond)
	commands := []Command{
		{LSN: 3, CommandType: 0, Args: []string{"name"}},
		{LSN: 2, CommandType: 1, Args: []string{"name", "Daniil"}},
		{LSN: 1, CommandType: 2, Args: []string{"name"}},
	}
	logsReader.EXPECT().ReadLogs().Return(commands, nil).Once()
	slices.SortFunc(commands, func(command1 Command, command2 Command) int {
		return command1.LSN - command2.LSN
	})

	res, err := wal.Recover()
	require.NoError(t, err)

	assert.Len(t, res, len(commands))
	for i := range res {
		assert.Equal(t, commands[i].Args, res[i].Args)
		assert.Equal(t, commands[i].CommandType, res[i].CommandType)
	}
	assert.Equal(t, wal.batch.lsn, commands[len(commands)-1].LSN+1)
}

func TestRecover_Error(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	wal, logsReader, _ := newTestWal(t, ctx, 10, 500*time.Millisecond)
	logsReader.EXPECT().ReadLogs().Return(nil, errors.New("recover error")).Once()

	_, err := wal.Recover()
	assert.Error(t, err)
}
