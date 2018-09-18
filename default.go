package task

import (
	"runtime"
	"time"
)

// Default task instance
var Default = NewTask(runtime.NumCPU() + 2)

// Add The specified time to execute
func Add(tim time.Time, task func()) (n *Node) {
	return Default.Add(tim, task)
}

// AddPeriodic Periodic execution
func AddPeriodic(perfunc func() time.Time, task func()) (n *Node) {
	return Default.AddPeriodic(perfunc, task)
}

// Cancel the task for this node
func Cancel(n *Node) {
	Default.Cancel(n)
}

// CancelAll Cancel all tasks
func CancelAll() {
	Default.CancelAll()
}

// Join Waiting for all tasks to finish
func Join() {
	Default.Join()
}

// Len returns the number of task
func Len() int {
	return Default.Len()
}

// List returns task list
func List() []*Node {
	return Default.List()
}

// Print task list
func Print() {
	Default.Print()
}
