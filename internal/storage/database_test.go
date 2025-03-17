package storage

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	engineerrors "github.com/DaniilZ77/InMemDB/internal/storage/engine"
	"github.com/DaniilZ77/InMemDB/internal/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute_Success(t *testing.T) {
	compute := mocks.NewCompute(t)
	engine := mocks.NewEngine(t)

	database, err := NewDatabase(compute, engine, slog.New(slog.NewJSONHandler(io.Discard, nil)))
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
				}, nil)
				engine.EXPECT().Get("name").Return("Daniil", nil).Once()
			},
		},
		{
			name:     "set command",
			command:  "set name Daniil",
			expected: "OK",
			mock: func() {
				compute.EXPECT().Parse("set name Daniil").Return(&parser.Command{
					Type: parser.SET,
					Args: []string{"name", "Daniil"},
				}, nil)
				engine.EXPECT().Set("name", "Daniil").Return().Once()
			},
		},
		{
			name:     "del command",
			command:  "del name",
			expected: "OK",
			mock: func() {
				compute.EXPECT().Parse("del name").Return(&parser.Command{
					Type: parser.DEL,
					Args: []string{"name"},
				}, nil)
				engine.EXPECT().Del("name").Return().Once()
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
				}, nil)
				engine.EXPECT().Get("name").Return("", engineerrors.ErrKeyNotFound).Once()
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
	compute := mocks.NewCompute(t)
	engine := mocks.NewEngine(t)

	database, err := NewDatabase(compute, engine, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	compute.EXPECT().Parse("get name").Return(nil, errors.New("internal error")).Once()

	res := database.Execute("get name")
	assert.Contains(t, res, "ERROR")
}
