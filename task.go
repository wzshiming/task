package task

import (
	"fmt"
	"time"

	"github.com/wzshiming/fork"
)

// TaskExit exit time
var TaskExit = time.Time{}

// unix 0
var unix0 = time.Unix(0, 0) // unix 0

// signal
var none = struct{}{}

// Task defined task sets
type Task struct {
	fork  *fork.Fork    // maximum thread control
	queue *list         // task queue
	curr  *Node         // currently waiting
	ins   chan struct{} // inserts signal
	iru   chan struct{} // is working
}

// NewTask create a new task that specifies the maximum number of fork
func NewTask(i int) *Task {
	if i < 1 {
		i = 1
	}
	i++
	return &Task{
		fork:  fork.NewFork(i),
		queue: newList(),
		ins:   make(chan struct{}, 1),
		iru:   make(chan struct{}, 1),
	}
}

// Join Waiting for all tasks to finish
func (t *Task) Join() {
	t.fork.Join()
}

// Cancel the task for this node
func (t *Task) Cancel(cn *Node) {
	n := &Node{}
	n.time = time.Unix(0, 0)
	n.task = func() {
		t.queue.Delete(cn)
	}
	t.add(n)
}

// CancelAll Cancel all tasks
func (t *Task) CancelAll() {
	n := &Node{}
	n.time = time.Unix(0, 0)
	n.task = func() {
		t.queue = newList()
	}
	t.add(n)
}

// add
func (t *Task) add(n *Node) *Node {
	now := time.Now()
	t.queue.InsertAndSort(n)

	if curr := t.curr; curr != nil && n.time.Before(curr.time) {
		t.flash()
	} else if n.time.Before(now) { // Refresh time if the first one is inserted
		t.flash()
	}

	select { // Start the administrative thread if it doesn't start
	case t.iru <- none:
		t.fork.Push(t.run)
	default:
	}
	return n
}

// Add The specified time to execute
func (t *Task) Add(tim time.Time, task func()) *Node {
	n := &Node{}
	n.time = tim
	n.task = func() {
		t.fork.Push(task)
	}
	n.name = time.Now().Format(time.RFC3339)
	return t.add(n)
}

// addPeriodic
func (t *Task) addPeriodic(perfunc func() time.Time, n *Node) *Node {
	if perfunc == nil {
		return nil
	}
	p := perfunc()
	if p.IsZero() {
		return nil
	}
	n.time = p
	return t.add(n)
}

// AddPeriodic Periodic execution
func (t *Task) AddPeriodic(perfunc func() time.Time, task func()) *Node {
	n := &Node{}
	n.task = func() {
		t.addPeriodic(perfunc, n)
		t.fork.Push(task)
	}
	n.name = time.Now().Format(time.RFC3339)
	return t.addPeriodic(perfunc, n)
}

// flash Reset the first task
func (t *Task) flash() {
	select {
	case t.ins <- none:
	default:
	}
}

// unflash
func (t *Task) unflash() {
	select {
	case <-t.ins:
	default:
	}
}

// run
func (t *Task) run() {
	timer := time.NewTimer(time.Hour)
	for {
		curr := t.queue.DeleteMin()
		if curr == nil { // End the thread if there are no tasks
			if t.Len() == 0 {
				break
			}
			continue
		}
		sub := curr.time.Sub(time.Now()) // Calculate sleep duration
		if sub <= 0 {                    // immediate
			curr.task()
			continue
		}
		timer.Reset(sub) // Reset timer
		t.curr = curr
		select {
		case <-t.ins: // A new task node is inserted
			t.curr = nil
			t.queue.InsertAndSort(curr)
		case <-timer.C: // Arrive at the recently executed task
			t.curr = nil
			curr.task()
			t.unflash()
		}
	}
	<-t.iru
}

// Len returns the number of task
func (t *Task) Len() int {
	b := 0
	if t.curr != nil {
		b++
	}
	return t.queue.Len() + b
}

// List returns task list
func (t *Task) List() []*Node {
	if t.curr != nil {
		return append([]*Node{t.curr}, t.queue.List()...)
	}
	return t.queue.List()
}

// First returns first on the list
func (t *Task) First() *Node {
	if t.curr != nil {
		return t.curr
	}
	return t.queue.Min()
}

// Last returns last on the list
func (t *Task) Last() *Node {
	return t.queue.Max()
}

// Print task list
func (t *Task) Print() {
	for _, v := range t.List() {
		fmt.Println(v)
	}
}
