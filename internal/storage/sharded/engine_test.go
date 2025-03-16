package sharded

import (
	"errors"
	"testing"

	"github.com/DaniilZ77/InMemDB/internal/storage/sharded/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLogShardsAmount = 3
)

func newBaseEngine(baseEngines []*mocks.BaseEngine) func() BaseEngine {
	i := 0
	return func() BaseEngine {
		oldI := i
		i++
		return baseEngines[oldI]
	}
}

func genBaseEngines(t *testing.T, size int) []*mocks.BaseEngine {
	engines := make([]*mocks.BaseEngine, 0, size)
	for range size {
		engines = append(engines, mocks.NewBaseEngine(t))
	}
	return engines
}

func TestGet_Success(t *testing.T) {
	t.Parallel()

	baseEngines := genBaseEngines(t, 1<<testLogShardsAmount)
	engine := NewShardedEngine(testLogShardsAmount, newBaseEngine(baseEngines))
	expected := "Daniil"

	baseEngines[engine.getHash("name")%(1<<testLogShardsAmount)].On("Get", "name").Return(expected, nil).Once()

	res, err := engine.Get("name")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, expected, res)
}

func TestGet_Fail(t *testing.T) {
	t.Parallel()

	baseEngines := genBaseEngines(t, 1<<testLogShardsAmount)
	engine := NewShardedEngine(testLogShardsAmount, newBaseEngine(baseEngines))

	baseEngines[engine.getHash("name")%(1<<testLogShardsAmount)].On("Get", "name").Return("", errors.New("internal error")).Once()

	_, err := engine.Get("name")
	assert.Error(t, err)
}

func TestSet(t *testing.T) {
	t.Parallel()

	baseEngines := genBaseEngines(t, 1<<testLogShardsAmount)
	engine := NewShardedEngine(testLogShardsAmount, newBaseEngine(baseEngines))

	baseEngines[engine.getHash("name")%(1<<testLogShardsAmount)].On("Set", "name", "Daniil").Return().Once()

	engine.Set("name", "Daniil")
}

func TestDel(t *testing.T) {
	t.Parallel()

	baseEngines := genBaseEngines(t, 1<<testLogShardsAmount)
	engine := NewShardedEngine(testLogShardsAmount, newBaseEngine(baseEngines))

	baseEngines[engine.getHash("name")%(1<<testLogShardsAmount)].On("Del", "name").Return().Once()

	engine.Del("name")
}
