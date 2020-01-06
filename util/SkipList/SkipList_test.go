package SkipList

import (
	"math/rand"
	"testing"
)

func TestSkipList(t *testing.T) {
	skiplist := SkipList(30, 5)
	t.Log(skiplist.Find(100))
	for i := 0; i < 20; i++ {
		t.Log(skiplist.Insert(rand.Float64() * 100))
	}
	sorted := skiplist.Traversal()
	t.Log(sorted)
	for _, node := range sorted {
		t.Log(node)
	}
}
