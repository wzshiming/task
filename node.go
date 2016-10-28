package task

import (
	"fmt"
	"llrb"
	"time"
)

// 任务节点
type node struct {
	time time.Time
	task func()
}

func (no *node) String() string {
	return fmt.Sprint(no.time)
}

// 执行的任务
func (no *node) Func() func() {
	return no.task
}

// 下次执行时间
func (no *node) Next() time.Time {
	return no.time
}

// 小于
func (no *node) Less(i llrb.Item) bool {
	switch i.(type) {
	case *node:
		b := i.(*node)
		return no != b && !no.time.After(b.time)
	default:
		return false
	}
}
