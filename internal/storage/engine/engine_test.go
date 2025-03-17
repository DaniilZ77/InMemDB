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

	engine.shards[engine.getHash("name")%(1<<testLogShardsAmount)].data.Store("name", "Daniil")

	res, err := engine.Get("name")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "Daniil", res)
}

func TestEngineGet_Fail(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(testLogShardsAmount)
	require.NoError(t, err)

	_, err = engine.Get("name")
	assert.Error(t, err)
}

func TesEngineSet(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(testLogShardsAmount)
	require.NoError(t, err)

	engine.Set("name", "Daniil")
	value, ok := engine.shards[engine.getHash("name")%(1<<testLogShardsAmount)].data.Load("name")
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}

func TestEngineDel(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(testLogShardsAmount)
	require.NoError(t, err)

	hash := engine.getHash("name") % (1 << testLogShardsAmount)

	engine.shards[hash].data.Store("name", "Daniil")

	engine.Del("name")
	_, ok := engine.shards[hash].data.Load("name")
	assert.False(t, ok)
}
