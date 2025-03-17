package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardGet_Success(t *testing.T) {
	t.Parallel()

	engine := NewShard()
	engine.data.Store("name", "Daniil")

	value, err := engine.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "Daniil", value)
}

func TestShardGet_Fail(t *testing.T) {
	t.Parallel()

	engine := NewShard()

	_, err := engine.Get("name")
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestShardSet(t *testing.T) {
	t.Parallel()

	engine := NewShard()

	engine.Set("name", "Daniil")

	value, ok := engine.data.Load("name")
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}

func TestShardDel(t *testing.T) {
	t.Parallel()

	engine := NewShard()
	engine.data.Store("name", "Daniil")

	engine.Del("name")

	_, ok := engine.data.Load("name")
	assert.False(t, ok)
}
