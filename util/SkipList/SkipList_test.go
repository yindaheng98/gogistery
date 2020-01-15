package SkipList

import (
	"math/rand"
	"testing"
)

func TestSkipList(t *testing.T) {
	skiplist := NewWithLevel(30, 5)
	t.Log(skiplist.Find(100))
	for i := 0; i < 20; i++ {
		t.Log(skiplist.Insert(rand.Float64() * 100))
	}
	sorted := skiplist.Traversal(16)
	t.Log(sorted)
	for _, node := range sorted {
		t.Log(node)
	}
	sorted = skiplist.TraversalAll()
	t.Log(sorted)
	for _, node := range sorted {
		t.Log(node)
	}
}
