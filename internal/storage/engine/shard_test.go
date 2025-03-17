package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardGet_Success(t *testing.T) {
	t.Parallel()

	engine := NewShard()
	engine.data["name"] = "Daniil"

	value, ok := engine.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}

func TestShardGet_NotFound(t *testing.T) {
	t.Parallel()

	engine := NewShard()

	_, err := engine.Get("name")
	assert.Equal(t, false, err)
}

func TestShardSet(t *testing.T) {
	t.Parallel()

	engine := NewShard()

	engine.Set("name", "Daniil")

	value, ok := engine.data["name"]
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}

func TestShardDel(t *testing.T) {
	t.Parallel()

	engine := NewShard()
	engine.data["name"] = "Daniil"

	engine.Del("name")

	_, ok := engine.data["name"]
	assert.False(t, ok)
}
