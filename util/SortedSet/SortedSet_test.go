package SortedSet

import (
	"fmt"
	"math/rand"
	"testing"
)

type testObj struct {
	data float64
}

func (o testObj) GetName() string {
	return fmt.Sprintf("I'm %.3f", o.data)
}

func TestSortedSet(t *testing.T) {
	zset := New(30)
	for i := 0; i < 20; i++ {
		e := new(testObj)
		e.data = rand.Float64()
		zset.Update(e, e.data)
	}
	var sorted = zset.Sorted(16)
	for _, e := range sorted {
		fmt.Println(e.GetName())
	}
	sorted = zset.SortedAll()
	for _, e := range sorted {
		fmt.Println(e.GetName())
	}
}
