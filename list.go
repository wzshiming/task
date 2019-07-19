package task

import (
	"sync"

	"github.com/wzshiming/llrb"
)

// list is an LLRB simulation of the table
type list struct {
	l   *llrb.LLRB
	mut sync.RWMutex
}

// newList is create a new list
func newList() *list {
	return &list{
		l: llrb.New(),
	}
}

// InsertAndSort is insert and sort
func (qu *list) InsertAndSort(n *Node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	qu.l.InsertNoReplace(n)
}

// DeleteMin is delete min node and returns it
func (qu *list) DeleteMin() (n *Node) {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	for qu.l.Len() != 0 && n == nil {
		n, _ = qu.l.DeleteMin().(*Node)
	}
	return n
}

// DeleteMin is delete the node and returns it
func (qu *list) Delete(n *Node) *Node {
	qu.mut.Lock()
	defer qu.mut.Unlock()
	b, _ := qu.l.Delete(n).(*Node)
	return b
}

// Min returns min node
func (qu *list) Min() *Node {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	b, _ := qu.l.Min().(*Node)
	return b
}

// Max returns max node
func (qu *list) Max() *Node {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	b, _ := qu.l.Max().(*Node)
	return b
}

// Len returns the list length
func (qu *list) Len() int {
	qu.mut.RLock()
	defer qu.mut.RUnlock()
	return qu.l.Len()
}

// List returns the list all node
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
