package task // import "gopkg.in/wzshiming/task.v2"

import (
	"fmt"
	"time"

	fork "gopkg.in/wzshiming/fork.v2"
)

var TaskExit = time.Time{} // exit time

var unix0 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local) // unix 0

var none = struct{}{} // signal

type Task struct {
	fork  *fork.Fork    // maximum thread control
	queue *list         // task queue
	curr  *Node         // currently waiting
	ins   chan struct{} // inserts signal
	iru   chan struct{} // is working
}

// NewTask
//  i: 线程数 最低为 1个线程
func NewTask(i int) *Task {
	if i < 1 {
		i = 1
	}
	i++
	return &Task{
		fork:  fork.NewFork(i),
		queue: NewList(),
		ins:   make(chan struct{}, 1),
		iru:   make(chan struct{}, 1),
	}
}

// Join Waiting for all tasks to finish
func (t *Task) Join() {
	t.fork.Join()
}

// Cancel
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
	t.flash()
	t.queue = NewList()
}

// add
func (t *Task) add(n *Node) *Node {

	t.queue.InsertAndSort(n) // 队列里插入

	if t.queue.Min() == n { // 如果插入到了第一个则刷新时间
		t.flash()
	}

	select { // 判断管理线程是否运行 如果没有则启动
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
		name: fmt.Sprint(task),
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
		name: fmt.Sprint(task),
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
		if t.curr == nil { // 如果没有任务了 结束线程
			if t.Len() == 0 {
				break
			}
			continue
		}
		sub := t.curr.time.Sub(time.Now()) // 计算 休眠时长
		if sub <= 0 {                      // 马上执行的
			t.curr.task()
			continue
		}
		timer.Reset(sub) // 重置定时器
		select {
		case <-t.ins: // 有新的 任务节点插入
			t.queue.InsertAndSort(t.curr)
		case <-timer.C: // 到达最近执行的任务
			t.curr.task()
			t.unflash()
		}
	}
	<-t.iru
}

// Len
func (t *Task) Len() int {
	b := 0
	if t.curr != nil {
		b++
	}
	return t.queue.Len() + b
}

// List
func (t *Task) List() []*Node {
	if t.curr != nil {
		return append([]*Node{t.curr}, t.queue.List()...)
	}
	return t.queue.List()
}

// Print
func (t *Task) Print() {
	for _, v := range t.List() {
		fmt.Printf(v.String())
	}
}
