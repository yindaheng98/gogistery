package SortedSet

import "gogistery/util/SkipList"

//一个用跳表和hashmap实现的有序集合
type SortedSet struct {
	skiplist     *SkipList.SkipList
	whosStringIs map[string]*SkipList.Node  //序列化Element->*node的map
	whosNodeIs   map[*SkipList.Node]Element //*node->*Element的map
}

func NewSortedSet(size uint64) *SortedSet {
	return &SortedSet{SkipList.NewSkipListWithC(size, 2),
		make(map[string]*SkipList.Node),
		make(map[*SkipList.Node]Element)}
}

//向集合中更新一个元素
func (set *SortedSet) Update(obj Element, weight float64) {
	str := obj.Stringify()
	set.remove(str)
	nodep := set.skiplist.Insert(weight)
	set.whosStringIs[str] = nodep
	set.whosNodeIs[nodep] = obj
}

//从集合中删除一个元素
func (set *SortedSet) Remove(obj Element) {
	set.remove(obj.Stringify())
}

func (set *SortedSet) remove(str string) {
	if nodep, ok := set.whosStringIs[str]; ok {
		set.skiplist.Delete(nodep)
		delete(set.whosStringIs, str)
		delete(set.whosNodeIs, nodep)
	}
}

func (set *SortedSet) GetWeight(obj Element) (float64, bool) {
	nodep, ok := set.whosStringIs[obj.Stringify()]
	if ok {
		return nodep.Data, true
	}
	return 0, false
}

func (set *SortedSet) Sorted(n uint64) []Element {
	return set.nodepsToElements(set.skiplist.Traversal(n))
}

func (set *SortedSet) SortedAll() []Element {
	return set.nodepsToElements(set.skiplist.TraversalAll())

}

func (set *SortedSet) nodepsToElements(nodeps []*SkipList.Node) []Element {
	length := len(nodeps)
	result := make([]Element, length)
	for i := 0; i < length; i++ {
		result[i] = set.whosNodeIs[nodeps[i]]
	}
	return result
}
