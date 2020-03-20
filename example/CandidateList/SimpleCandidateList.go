package CandidateList

import (
	"context"
	"github.com/yindaheng98/go-utility/SortedSet"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//SimpleCandidateList is a simple implementation of RegistryCandidateList.
//This implementation sort registries by a weight.
type SimpleCandidateList struct {
	DefaultTimeout time.Duration
	DefaultRetryN  uint64
	set            chan *SortedSet.SortedSet
	waitGroup      chan bool
}

//NewSimpleCandidateList returns the pointer to a SimpleCandidateList, with initRegistry in it.
func NewSimpleCandidateList(size uint64, initRegistry protocol.RegistryInfo) *SimpleCandidateList {
	list := NewEmptySimpleCandidateList(size)
	list.StoreCandidates(context.Background(), []protocol.RegistryInfo{initRegistry})
	return list
}

//NewEmptySimpleCandidateList returns the pointer to a empty SimpleCandidateList.
func NewEmptySimpleCandidateList(size uint64) *SimpleCandidateList {
	list := &SimpleCandidateList{
		DefaultTimeout: 1e9,
		DefaultRetryN:  10,
		set:            make(chan *SortedSet.SortedSet, 1),
		waitGroup:      make(chan bool),
	}
	set := SortedSet.New(size)
	list.set <- set
	return list
}

//StoreCandidates adds or updates registries in the list candidates and change the registries' weight
//according to the following principle:
//
//1. If a candidate registry is not exist, add the registry and give the set weight to 0
//
//2. If a candidate registry is exist, update the registry but do not change its weight
//
//3. When added or updated, all the other candidate registries' weight increase by 1
func (list *SimpleCandidateList) StoreCandidates(ctx context.Context, candidates []protocol.RegistryInfo) {
	set := <-list.set
	defer func() {
		list.set <- set                  //完成后放回队列
		close(list.waitGroup)            //唤醒所有等待的获取队列进程
		list.waitGroup = make(chan bool) //然后再阻塞之
	}()
	set.DeltaUpdateAll(-1) //先让所有元素优先级下降1
	for _, candidate := range candidates {
		el := element{candidate}
		if _, exists := set.GetWeight(el); exists {
			set.DeltaUpdate(el, 1) //然后让响应中给出的元素优先级上升1
		} else {
			set.Update(el, 0) //或者加入新元素
		}
	}
}

//GetCandidate select the candidate registry whose weight is smallest and is not in "excepts" as result and delete it from SimpleCandidateList.
//If there is no candidate meet the above conditions, block until a eligible candidate added
func (list *SimpleCandidateList) GetCandidate(ctx context.Context, excepts []protocol.RegistryInfo) (candidate protocol.RegistryInfo, initTimeout time.Duration, initRetryN uint64) {
	for {
		set := <-list.set
		for _, except := range excepts {
			set.Remove(element{except}) //先移除所有除外元素
		}
		if els := set.Sorted(1); len(els) > 0 { //然后看是否还有剩余元素
			el := els[0].(element) //有就作为返回结果
			set.Remove(el)         //然后删除返回结果
			list.set <- set        //并将集合放回队列
			return el.RegistryInfo, list.DefaultTimeout, list.DefaultRetryN
		}
		list.set <- set  //并将集合放回队列
		<-list.waitGroup //等待更新
	}
}

//Delete a candidate from CandidateList. Implemention of CandidateList.DeleteCandidate.
func (list *SimpleCandidateList) DeleteCandidate(ctx context.Context, info protocol.RegistryInfo) {
	var set *SortedSet.SortedSet
	select {
	case <-ctx.Done(): //若要结束
		return //则直接退出
	case set = <-list.set: //取出集合
		defer func() { list.set <- set }() //结束时放回
	}
	set.Remove(element{info}) //删除元素
}
