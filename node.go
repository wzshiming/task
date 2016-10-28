package task

import (
	"fmt"
	"llrb"
	"time"
)

// 任务节点
type Node struct {
	time time.Time
	task func()
	name string
}

func (no *Node) String() string {
	return fmt.Sprint(no.time)
}

// 执行的任务
func (no *Node) Func() func() {
	return no.task
}

// 下次执行时间
func (no *Node) Next() time.Time {
	return no.time
}

// 设置名字
func (no *Node) SetName(name string) {
	no.name = name
}

// 名字
func (no *Node) Name() string {
	return no.name
}

// 小于
func (no *Node) Less(i llrb.Item) bool {
	switch i.(type) {
	case *Node:
		b := i.(*Node)
		return no != b && !no.time.After(b.time)
	default:
		return false
	}
}
