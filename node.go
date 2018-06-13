package task

import (
	"fmt"
	"time"

	"github.com/wzshiming/llrb"
)

// Node is store tasks
type Node struct {
	time time.Time
	task func()
	name string
}

// String returns strings
func (no *Node) String() string {
	return fmt.Sprintf("%v %v", no.time, no.name)
}

// Func returns tasks function
func (no *Node) Func() func() {
	return no.task
}

// Next returns next time
func (no *Node) Next() time.Time {
	return no.time
}

// SetName sets node name
func (no *Node) SetName(name string) {
	no.name = name
}

// Name returns the node name
func (no *Node) Name() string {
	return no.name
}

// Less returns compare the time with another node
func (no *Node) Less(i llrb.Item) bool {
	switch b := i.(type) {
	case *Node:
		return no != b && !no.time.After(b.time)
	default:
		return false
	}
}
