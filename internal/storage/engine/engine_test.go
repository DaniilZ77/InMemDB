package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLogShardsAmount = 3
)

func TestEngineGet_Success(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(testLogShardsAmount)
	require.NoError(t, err)

	value1, value2 := "Daniil", "Ilia"
	engine.shards[engine.getHash("name")%(1<<testLogShardsAmount)].Set(1, "name", &value1)
	engine.shards[engine.getHash("name")%(1<<testLogShardsAmount)].Set(2, "name", &value2)

	res, ok := engine.Get(3, "name")
	require.True(t, ok)
	require.NotNil(t, res)
	assert.Equal(t, "Ilia", res)
}

func TestEngineGet_NotFound(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(testLogShardsAmount)
	require.NoError(t, err)

	_, ok := engine.Get(1, "name")
	assert.False(t, ok)
}

func TestEngineSet(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(testLogShardsAmount)
	require.NoError(t, err)

	value := "Daniil"
	engine.Set(1, "name", &value)
	value, ok := engine.shards[engine.getHash("name")%(1<<testLogShardsAmount)].Get(2, "name")
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}
