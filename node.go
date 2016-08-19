package task

import (
	"fmt"
	"llrb"
	"time"
	"unsafe"
)

// 任务节点
type node struct {
	time   time.Time
	prefix func()
	task   func()
}

func (no *node) String() string {
	return fmt.Sprint(no.time, no.task)
}

// 小于
func (no *node) Less(i llrb.Item) bool {
	switch i.(type) {
	case *node:
		b := i.(*node)
		if uintptr(unsafe.Pointer(no)) == uintptr(unsafe.Pointer(b)) {
			return false
		}
		return !no.time.After(b.time)
	default:
		return false
	}
}
