package SkipList

type Node struct {
	Data float64
	prev []*Node
	next []*Node
}

func NewNode(Data float64, level uint64) *Node {
	if level < 1 {
		level = 1
	}
	return &Node{Data, make([]*Node, level), make([]*Node, level)}
}
