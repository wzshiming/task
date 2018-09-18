package task

import (
	"fmt"
	"time"

	"github.com/wzshiming/fork"
)

// TaskExit exit time
var TaskExit = time.Time{}

// unix 0
var unix0 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local) // unix 0

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
func (t *Task) Cancel(n *Node) {
	t.add(&Node{
		time: time.Unix(0, 0),
		task: func() {
			t.queue.Delete(n)
		},
	})
}

// CancelAll Cancel all tasks
func (t *Task) CancelAll() {
	t.add(&Node{
		time: time.Unix(0, 0),
		task: func() {
			t.queue = newList()
		},
	})
}

// add
func (t *Task) add(n *Node) *Node {

	t.queue.InsertAndSort(n)

	if t.queue.Min() == n && n.time.Before(time.Now()) { // Refresh time if the first one is inserted
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
	return t.add(&Node{
		time: tim,
		task: func() {
			t.fork.Push(task)
		},
		name: time.Now().Format(time.RFC3339),
	})
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
func (t *Task) AddPeriodic(perfunc func() time.Time, task func()) (n *Node) {
	n = &Node{
		task: func() {
			t.addPeriodic(perfunc, n)
			t.fork.Push(task)
		},
		name: time.Now().Format(time.RFC3339),
	}
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
		t.curr = t.queue.DeleteMin()
		if t.curr == nil { // End the thread if there are no tasks
			if t.Len() == 0 {
				break
			}
			continue
		}
		sub := t.curr.time.Sub(time.Now()) // Calculate sleep duration
		if sub <= 0 {                      // immediate
			t.curr.task()
			continue
		}
		timer.Reset(sub) // Reset timer
		select {
		case <-t.ins: // A new task node is inserted
			t.queue.InsertAndSort(t.curr)
		case <-timer.C: // Arrive at the recently executed task
			t.curr.task()
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

// Print task list
func (t *Task) Print() {
	for _, v := range t.List() {
		fmt.Println(v)
	}
}
