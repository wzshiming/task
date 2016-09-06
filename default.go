package task

import (
	"runtime"
	"time"
)

// 全局默认任务管理
var Default = NewTask(runtime.NumCPU() + 2)

// 添加单次定时任务
func Add(tim time.Time, task func()) (n *node) {
	return Default.Add(tim, task)
}

// 添加间隔时间任务
func AddPeriodic(perfunc func() time.Time, task func()) (n *node) {
	return Default.AddPeriodic(perfunc, task)
}

// 取消任务
func Cancel(n *node) {
	Default.Cancel(n)
}

// 等待任务结束
func Join() {
	Default.Join()
}

func Len() int {
	return Default.Len()
}
