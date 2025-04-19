package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardGet_Success(t *testing.T) {
	t.Parallel()

	value := "Daniil"
	engine := NewShard()
	engine.data.Insert(VersionedKey{
		key:     "name",
		version: 1,
	}, &value)

	value, ok := engine.Get(2, "name")
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}

func TestShardGet_NotFound(t *testing.T) {
	t.Parallel()

	engine := NewShard()

	_, err := engine.Get(1, "name")
	assert.Equal(t, false, err)
}

func TestShardSet(t *testing.T) {
	t.Parallel()

	engine := NewShard()

	value := "Daniil"
	engine.Set(1, "name", &value)

	value, ok := engine.data.Find(VersionedKey{
		key:     "name",
		version: 2,
	})
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
}
