package parser

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Success(t *testing.T) {
	t.Parallel()

	p, err := NewParser(slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	tests := []struct {
		name     string
		command  string
		expected *Command
	}{
		{
			name:    "get command",
			command: "get name",
			expected: &Command{
				Type: GET,
				Args: []string{"name"},
			},
		},
		{
			name:    "set command",
			command: "set name Daniil",
			expected: &Command{
				Type: SET,
				Args: []string{"name", "Daniil"},
			},
		},
		{
			name:    "delete command",
			command: "del name",
			expected: &Command{
				Type: DEL,
				Args: []string{"name"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := p.Parse(tt.command)
			if assert.NoError(t, err) {
				assert.NotNil(t, cmd)
				assert.Equal(t, *tt.expected, *cmd)
			}
		})
	}
}

func TestParse_Fail(t *testing.T) {
	t.Parallel()

	p, err := NewParser(slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "empty command",
			command: "",
		},
		{
			name:    "bad command type",
			command: "sat name Daniil",
		},
		{
			name:    "bad amount of args",
			command: "get name Daniil",
		},
		{
			name:    "bad amount of args",
			command: "set name",
		},
		{
			name:    "bad amount of args",
			command: "del",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := p.Parse(tt.command)
			assert.ErrorIs(t, err, ErrInvalidCommand)
		})
	}
}
