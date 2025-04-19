package engine

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toSkipList(keys []VersionedKey, values []string) *SkipList {
	sl := &SkipList{head: &SkipNode{}}
	cur := sl.head
	for i := range keys {
		cur.next = &SkipNode{key: keys[i], value: &values[i]}
		cur = cur.next
	}
	return sl
}

func TestCreateSkipList(t *testing.T) {
	t.Parallel()

	keys := []VersionedKey{{"key1", 1}, {"key1", 2}, {"key1", 3}, {"key2", 1}, {"key3", 1}, {"key4", 1}, {"key5", 1}}
	values := []string{"value11", "value12", "value13", "value2", "value3", "value4", "value5"}
	sl := toSkipList(keys, values)
	sl = NewSkipList(sl)

	for i := range keys {
		key := keys[i]
		key.version++
		value, ok := sl.Find(key)
		assert.True(t, ok)
		assert.Equal(t, values[i], value)
	}

	for range 100 {
		_, ok := sl.Find(VersionedKey{strconv.Itoa(rand.Int()), 1})
		assert.False(t, ok)
	}
}

func TestInsertSkipList(t *testing.T) {
	t.Parallel()

	keys := []VersionedKey{{"key1", 1}, {"key1", 2}, {"key1", 3}, {"key2", 1}, {"key3", 1}, {"key4", 1}, {"key5", 1}}
	values := []string{"value11", "value12", "value13", "value2", "value3", "value4", "value5"}
	sl := toSkipList(keys, values)
	sl = NewSkipList(sl)

	for i := len(keys); i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		sl.Insert(VersionedKey{key, 1}, &value)
		foundValue, ok := sl.Find(VersionedKey{key, 2})
		assert.True(t, ok)
		assert.Equal(t, value, foundValue)
	}
}

func TestSkipList_InsertEqual(t *testing.T) {
	t.Parallel()

	sl := NewSkipList(nil)
	var (
		value1 = "value1"
		value2 = "value2"
		value3 = "value3"
		value4 = "value4"
	)
	sl.Insert(VersionedKey{"key1", 1}, &value1)
	sl.Insert(VersionedKey{"key1", 1}, &value2)
	sl.Insert(VersionedKey{"key1", 1}, &value3)
	sl.Insert(VersionedKey{"key1", 1}, &value4)
	value, ok := sl.Find(VersionedKey{"key1", 2})
	assert.True(t, ok)
	assert.Equal(t, "value4", value)
}

func TestSkipList_ExistsBetween(t *testing.T) {
	t.Parallel()

	sl := NewSkipList(nil)
	var (
		value1 = "value1"
		value2 = "value2"
		value3 = "value3"
	)
	sl.Insert(VersionedKey{"key1", 1}, &value1)
	sl.Insert(VersionedKey{"key1", 3}, &value2)
	sl.Insert(VersionedKey{"key1", 4}, &value3)

	assert.True(t, sl.ExistsBetween(VersionedKey{"key1", 1}, VersionedKey{"key1", 3}))
	assert.False(t, sl.ExistsBetween(VersionedKey{"key1", 2}, VersionedKey{"key1", 3}))
}
