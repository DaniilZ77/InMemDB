package engine

import (
	"math/rand/v2"
	"strings"
)

type VersionedKey struct {
	key     string
	version int64
}

type SkipNode struct {
	down  *SkipNode
	next  *SkipNode
	key   VersionedKey
	value *string
}

type SkipList struct {
	head *SkipNode
}

func cmp(key1 VersionedKey, key2 VersionedKey) int {
	if key1.key == key2.key {
		return int(key1.version - key2.version)
	}
	return strings.Compare(key1.key, key2.key)
}

func buildLevel(sl *SkipList) *SkipList {
	nextLevel := &SkipList{head: &SkipNode{down: sl.head}}
	i := sl.head.next.next
	cur := nextLevel.head
	for i != nil && i.next != nil {
		cur.next = &SkipNode{next: cur.next, down: i, key: i.key}
		cur = cur.next
		i = i.next.next
	}
	return nextLevel
}

func NewSkipList(sl *SkipList) *SkipList {
	if sl == nil || sl.head == nil {
		return &SkipList{head: &SkipNode{}}
	}

	for sl.head.next != nil {
		sl = buildLevel(sl)
	}
	return sl
}

func (sl *SkipList) Insert(key VersionedKey, value *string) {
	sn := sl.head.find(key).next
	if sn != nil && cmp(sn.key, key) == 0 {
		sn.value = value
		return
	}

	inserted := sl.head.insert(key, value)
	if inserted != nil {
		*sl = SkipList{head: &SkipNode{next: inserted, down: sl.head}}
	}
}

func (sn *SkipNode) insert(key VersionedKey, value *string) *SkipNode {
	for sn.next != nil && cmp(sn.next.key, key) < 0 {
		sn = sn.next
	}
	var inserted *SkipNode
	if sn.down == nil {
		sn.next = &SkipNode{next: sn.next, key: key, value: value}
		inserted = sn.next
	} else {
		inserted = sn.down.insert(key, value)
		if inserted != nil {
			sn.next = &SkipNode{next: sn.next, down: inserted, key: inserted.key}
		}
	}
	if inserted != nil && rand.Float64() < 0.5 {
		return inserted
	}
	return nil
}

func (sl *SkipList) Delete(key VersionedKey) {
	sl.head.delete(key)
}

func (sn *SkipNode) delete(key VersionedKey) {
	panic("unimplemented")
}

func (sl *SkipList) Find(key VersionedKey) (string, bool) {
	sn := sl.head.find(key)
	if sn.value == nil {
		return "", false
	}
	return *sn.value, sn.key.key == key.key
}

func (sn *SkipNode) find(key VersionedKey) *SkipNode {
	for sn.next != nil && cmp(sn.next.key, key) < 0 {
		sn = sn.next
	}
	if sn.down == nil {
		return sn
	}

	return sn.down.find(key)
}

func (sl *SkipList) ExistsBetween(key1 VersionedKey, key2 VersionedKey) bool {
	lower := sl.head.find(key1)
	upper := sl.head.find(key2)
	return lower != upper
}
