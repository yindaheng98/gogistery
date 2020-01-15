package SkipList

import (
	"math"
	"time"
)

type SkipList struct {
	root      *Node      //根节点指针
	n         uint64     //节点总数
	level     uint64     //预估的索引层数
	randLevel *RandLevel //索引层数生成器
}

//构造一个跳表
//
//listSize是预计将在的跳表中存入的节点总数，indexLevel是索引的最大层数（总层数=索引层数+1）
func NewWithLevel(listSize, indexLevel uint64) *SkipList {
	C := uint64(math.Ceil(math.Pow(float64(listSize), 1.0/float64(indexLevel))))
	return &SkipList{nil, 0, indexLevel + 1,
		NewRandomLevel(C, indexLevel, time.Now().UnixNano())}
}

//构造一个跳表
//
//listSize是预计将在的跳表中存入的节点总数，C是索引的衰减系数（下一层索引数=上一层索引数/C）
func NewWithC(listSize, C uint64) *SkipList {
	indexLevel := uint64(math.Ceil(math.Log(float64(listSize)) / math.Log(float64(C))))
	return &SkipList{nil, 0, indexLevel + 1,
		NewRandomLevel(C, indexLevel, time.Now().UnixNano())}
}

//找到各层index中大小小于data的最大节点的指针
func (sl *SkipList) Find(data float64) *Node {
	result := sl.find(data)
	if result == nil || len(result) < 1 {
		return nil
	}
	return result[0]
}

//找到各层index中大小小于data的最大节点的指针
func (sl *SkipList) find(data float64) []*Node {
	if sl.root == nil { //如果链表为空
		return nil //则直接返回空
	}

	//链表不为空才能开始初始化
	level := len(sl.root.next)     //根节点索引层数即时最大索引层数
	result := make([]*Node, level) //初始化结果index表
	if data < sl.root.Data {       //如果链表中没有这样的节点就直接返回
		return result
	}

	//有这样的节点才能开始查找
	p := sl.root     //初始化当前指针
	pLevel := level  //初始化当前指针所在层数
	for pLevel > 0 { //循环直到pLevel到了第0层
		pLevel -= 1                           //index向下走一层
		next := p.next[pLevel]                //初始化该层的下一个节点指针
		for next != nil && next.Data < data { //如果后面有节点并且其值比data小
			p = next //就往后走一步
			next = p.next[pLevel]
		} //走到头了就退出，此时的p即第pLevel层要找的节点指针
		result[pLevel] = p //记录这个指针
	}
	return result
}

//插入一个数据
func (sl *SkipList) Insert(data float64) *Node {
	sl.n++
	pres := sl.find(data)      //查找插入点
	presN := uint64(len(pres)) //插入节点的数量

	if pres == nil { //查找返回了空，说明链表为空
		sl.root = NewNode(data, sl.level) //那就直接给root赋值
		return sl.root
	}

	//链表不为空才能开始正常赋值
	level := sl.randLevel.Rand() + 1
	result := NewNode(data, level) //要返回的值（永远返回新指针）
	insert := result               //要执行插入操作的值

	//返回的第一个指针就为空
	if pres[0] == nil { //说明要在根节点前插
		insert.prev = sl.root.prev
		insert.next = sl.root.next //首先复制根节点的前后指针
		sl.root.prev = make([]*Node, level)
		sl.root.next = make([]*Node, level)  //然后新建根节点的前后指针
		insert = sl.root                     //“偷梁换柱”：把根节点值提出来作为要插入的值
		sl.root = result                     //然后将根节点值替换为新值
		for i := uint64(0); i < presN; i++ { //然后更新前置节点表
			pres[i] = sl.root
		}
	}

	//最后执行插入操作
	for i := uint64(0); i < level; i++ {
		insert.prev[i] = pres[i]
		insert.next[i] = pres[i].next[i]
		pres[i].next[i] = insert
	}
	return result
}

//升序遍历
func (sl *SkipList) Traversal(n uint64) []*Node {
	if sl.n < n {
		n = sl.n
	}
	result := make([]*Node, n)
	node := sl.root
	for i := uint64(0); i < n && node != nil; i++ {
		result[i] = node
		node = node.next[0]
	}
	return result
}

func (sl *SkipList) TraversalAll() []*Node {
	return sl.Traversal(sl.n)
}

func (sl *SkipList) Delete(node *Node) {
	prev := node.prev
	next := node.next
	length := len(prev)
	for i := 0; i < length; i++ {
		if prev[i] != nil {
			prev[i].next[i] = next[i]
		}
		if next[i] != nil {
			next[i].prev[i] = prev[i]
		}
	}
	if node == sl.root {
		sl.root = node.next[0]
	}
}
