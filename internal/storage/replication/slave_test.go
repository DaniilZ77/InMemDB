package replication

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	t.Parallel()

	disk, client := NewMockSegmentManager(t), NewMockClient(t)

	disk.EXPECT().LastSegment().Return("segment.log", nil).Once()
	client.EXPECT().Close().Return(nil).Once()

	slave, err := NewSlave(50*time.Millisecond, client, disk, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	expectedCommands := []wal.Command{
		{LSN: 1, CommandType: 1, Args: []string{"name", "Daniil"}},
		{LSN: 2, CommandType: 0, Args: []string{"name"}},
		{LSN: 3, CommandType: 2, Args: []string{"name"}},
	}
	encodedExpectedCommands, err := common.Encode(expectedCommands)
	require.NoError(t, err)

	request, err := common.Encode(NewRequest("segment.log"))
	require.NoError(t, err)
	response, err := common.Encode(NewSuccessResponse("segment.log", encodedExpectedCommands))
	require.NoError(t, err)

	client.EXPECT().Send(request).Return(response, nil)
	disk.EXPECT().WriteFile("segment.log", encodedExpectedCommands).Return(nil)

	wg := sync.WaitGroup{}
	wg.Add(2)

	commandsNumber := 0

	go func() {
		defer wg.Done()
		for commands := range slave.GetReplicationStream() {
			commandsNumber++
			assert.Equal(t, expectedCommands, commands)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go func() {
		defer wg.Done()
		err = slave.Start(ctx)
		assert.NoError(t, err)
	}()

	wg.Wait()
	assert.LessOrEqual(t, commandsNumber, 10)
	assert.Greater(t, commandsNumber, 0)
}
