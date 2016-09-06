package task

import (
	"runtime"
	"time"
)

var Default = NewTask(runtime.NumCPU() + 2)

func Add(tim time.Time, task func()) (n *node) {
	return Default.Add(tim, task)
}

func AddPeriodic(perfunc func() time.Time, task func()) (n *node) {
	return Default.AddPeriodic(perfunc, task)
}

func Cancel(n *node) {
	Default.Cancel(n)
}

func Join() {
	Default.Join()
}

func Len() int {
	return Default.Len()
}
