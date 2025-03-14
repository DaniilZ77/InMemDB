package baseengine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet_Success(t *testing.T) {
	t.Parallel()

	engine := NewEngine()
	engine.data["name"] = "Daniil"

	value, err := engine.Get("name")
	if assert.NoError(t, err) {
		assert.Equal(t, "Daniil", *value)
	}
}

func TestGet_Fail(t *testing.T) {
	t.Parallel()

	engine := NewEngine()

	_, err := engine.Get("name")
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestSet(t *testing.T) {
	t.Parallel()

	engine := NewEngine()

	engine.Set("name", "Daniil")

	value, ok := engine.data["name"]
	if assert.True(t, ok) {
		assert.Equal(t, "Daniil", value)
	}
}

func TestDel(t *testing.T) {
	t.Parallel()

	engine := NewEngine()
	engine.data["name"] = "Daniil"

	engine.Del("name")

	_, ok := engine.data["name"]
	assert.False(t, ok)
}
