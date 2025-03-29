package wal

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogsManager_Success(t *testing.T) {
	t.Parallel()

	mockDisk := NewMockDisk(t)
	logsManager := NewLogsManager(mockDisk, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	commands := []Command{
		{LSN: 1, CommandType: 0, Args: []string{"name"}},
		{LSN: 2, CommandType: 1, Args: []string{"name", "Daniil"}},
		{LSN: 3, CommandType: 2, Args: []string{"name"}},
	}

	var encodedCommands []byte
	mockDisk.EXPECT().WriteSegment(mock.MatchedBy(func(data []byte) bool {
		encodedCommands = data
		return true
	})).Return(nil).Once()

	err := logsManager.WriteLogs(commands)
	require.NoError(t, err)

	mockDisk.EXPECT().ReadSegments().Return(encodedCommands, nil).Once()

	decodedCommands, err := logsManager.ReadLogs()
	require.NoError(t, err)
	assert.Equal(t, commands, decodedCommands)
}

func TestLogsManager_Error(t *testing.T) {
	mockDisk := NewMockDisk(t)
	logsManager := NewLogsManager(mockDisk, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	expectedErr := errors.New("logs manager error")

	tests := []struct {
		name string
		call func() error
		mock func()
	}{
		{
			name: "read error",
			call: func() error {
				_, err := logsManager.ReadLogs()
				return err
			},
			mock: func() {
				mockDisk.EXPECT().ReadSegments().Return(nil, expectedErr).Once()
			},
		},
		{
			name: "write error",
			call: func() error {
				return logsManager.WriteLogs([]Command{})
			},
			mock: func() {
				mockDisk.EXPECT().WriteSegment(mock.Anything).Return(expectedErr).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()
			err := tt.call()
			assert.ErrorIs(t, expectedErr, err)
		})
	}
}
