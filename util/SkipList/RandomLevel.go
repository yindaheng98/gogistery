package SkipList

import (
	"math"
	"math/rand"
)

type RandLevel struct {
	logC    float64
	min     float64
	randSrc rand.Source
}

//新建一个随机层数生成器，算法详情参见https://yindaheng98.github.io/%E6%95%B0%E5%AD%A6/SkipList.html
//
//下一层的索引数量是上一层的1/C，共Level层
func NewRandomLevel(C, Level uint64, seed int64) *RandLevel {
	if C <= 1 {
		C = 2
	}
	c := float64(C)
	level := float64(Level)
	return &RandLevel{math.Log(c), 1.0 / math.Pow(c, level+1), rand.NewSource(seed)}
}

func (rl *RandLevel) Rand() uint64 {
	X0 := float64(0)
	for X0 == 0 {
		X0 = rand.New(rl.randSrc).Float64()
	}
	X0 = X0*(1-rl.min) + rl.min
	X := math.Floor(-math.Log(X0) / rl.logC)
	return uint64(X)
}
