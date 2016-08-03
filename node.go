package task

import (
	"llrb"
	"time"
	"unsafe"
)

// 任务节点
type node struct {
	tim time.Time
	fun func()
}

// 小于
func (no *node) Less(i llrb.Item) bool {
	switch i.(type) {
	case *node:
		b := i.(*node)
		if uintptr(unsafe.Pointer(no)) == uintptr(unsafe.Pointer(b)) {
			return false
		}
		return no.tim.Before(b.tim) || no.tim.Equal(b.tim)
	default:
		return false
	}
}
