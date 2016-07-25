package task

import (
	"llrb"
	"sync"
)

type Queue struct {
	l   *llrb.LLRB
	mut sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		l: llrb.New(),
	}
}

func (qu *Queue) InsertAndSort(n *node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	qu.l.InsertNoReplace(n)
}

func (qu *Queue) DeleteMin() *node {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	b, _ := qu.l.DeleteMin().(*node)
	return b
}

func (qu *Queue) Min() *node {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	b, _ := qu.l.Min().(*node)
	return b
}

func (qu *Queue) Len() int {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	return qu.l.Len()
}
