package task

import (
	"llrb"
	"sync"
)

// 任务队列
type List struct {
	l   *llrb.LLRB
	mut sync.RWMutex
}

// 新的任务队列
func NewList() *List {
	return &List{
		l: llrb.New(),
	}
}

// 插入并排序
func (qu *List) InsertAndSort(n *node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	qu.l.InsertNoReplace(n)
}

// 删除最小的
func (qu *List) DeleteMin() (n *node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	for qu.l.Len() != 0 && n == nil {
		n, _ = qu.l.DeleteMin().(*node)
	}
	return n
}

// 删除某个节点
func (qu *List) Delete(n *node) *node {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	b, _ := qu.l.Delete(n).(*node)
	return b
}

// 获取最小的
func (qu *List) Min() *node {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	b, _ := qu.l.Min().(*node)
	return b
}

// 长度
func (qu *List) Len() int {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	return qu.l.Len()
}
