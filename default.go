package task

import (
	"runtime"
	"time"
)

// Default
var Default = NewTask(runtime.NumCPU() + 2)

// Add The specified time to execute
func Add(tim time.Time, task func()) (n *Node) {
	return Default.Add(tim, task)
}

// AddPeriodic Periodic execution
func AddPeriodic(perfunc func() time.Time, task func()) (n *Node) {
	return Default.AddPeriodic(perfunc, task)
}

// Cancel
func Cancel(n *Node) {
	Default.Cancel(n)
}

// Join Waiting for all tasks to finish
func Join() {
	Default.Join()
}

// Len
func Len() int {
	return Default.Len()
}

// List
func List() []*Node {
	return Default.List()
}

// Print
func Print() error {
	return Default.Print()
}
