package task

import (
	"sync"

	"github.com/wzshiming/llrb"
)

// 任务队列
type list struct {
	l   *llrb.LLRB
	mut sync.RWMutex
}

// 新的任务队列
func NewList() *list {
	return &list{
		l: llrb.New(),
	}
}

// 插入并排序
func (qu *list) InsertAndSort(n *Node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	qu.l.InsertNoReplace(n)
}

// 删除最小的
func (qu *list) DeleteMin() (n *Node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	for qu.l.Len() != 0 && n == nil {
		n, _ = qu.l.DeleteMin().(*Node)
	}
	return n
}

// 删除某个节点
func (qu *list) Delete(n *Node) *Node {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	b, _ := qu.l.Delete(n).(*Node)
	return b
}

// 获取最小的
func (qu *list) Min() *Node {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	b, _ := qu.l.Min().(*Node)
	return b
}

// 长度
func (qu *list) Len() int {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	return qu.l.Len()
}

// 节点列表
func (qu *list) List() []*Node {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	ns := make([]*Node, 0, qu.l.Len())
	qu.l.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		b, _ := i.(*Node)
		if b != nil {
			ns = append(ns, b)
		}
		return true
	})
	return ns
}
