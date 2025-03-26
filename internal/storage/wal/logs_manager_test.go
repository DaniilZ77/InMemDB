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

	disk := NewMockDisk(t)
	logsManager := NewLogsManager(disk, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	commands := []Command{
		{LSN: 1, CommandType: 0, Args: []string{"name"}},
		{LSN: 2, CommandType: 1, Args: []string{"name", "Daniil"}},
		{LSN: 3, CommandType: 2, Args: []string{"name"}},
	}

	var encodedCommands []byte
	disk.EXPECT().Write(mock.MatchedBy(func(data []byte) bool {
		encodedCommands = data
		return true
	})).Return(nil).Once()

	err := logsManager.Write(commands)
	require.NoError(t, err)

	disk.EXPECT().Read().Return(encodedCommands, nil).Once()

	decodedCommands, err := logsManager.Read()
	require.NoError(t, err)
	assert.Equal(t, commands, decodedCommands)
}

func TestLogsManager_Error(t *testing.T) {
	disk := NewMockDisk(t)
	logsManager := NewLogsManager(disk, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	expectedErr := errors.New("logs manager error")

	tests := []struct {
		name string
		call func() error
		mock func()
	}{
		{
			name: "read error",
			call: func() error {
				_, err := logsManager.Read()
				return err
			},
			mock: func() {
				disk.EXPECT().Read().Return(nil, expectedErr).Once()
			},
		},
		{
			name: "write error",
			call: func() error {
				return logsManager.Write([]Command{})
			},
			mock: func() {
				disk.EXPECT().Write([]byte(nil)).Return(expectedErr).Once()
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
