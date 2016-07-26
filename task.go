package task

import (
	"time"

	"github.com/wzshiming/fork"
)

type Task struct {
	queue *Queue
	cha   chan struct{}
	f     *fork.Fork
}

func NewTask(i int) *Task {
	t := &Task{
		queue: NewQueue(),
		cha:   make(chan struct{}, 1),
		f:     fork.NewFork(i),
	}
	go t.run()
	return t
}

func (t *Task) Add(tim time.Time, f func()) {
	n := &node{
		Time: tim,
		fun:  f,
	}

	t.queue.InsertAndSort(n)
	if t.queue.Min() == n {
		select {
		case t.cha <- struct{}{}:
		default:
		}
	}
}

func (t *Task) AddPeriodic(perfunc func() time.Time, f func()) {
	p := perfunc()
	if p.IsZero() {
		return
	}
	t.Add(p, func() {
		f()
		t.AddPeriodic(perfunc, f)
	})
}

func (t *Task) sleep() {
	t.Add(time.Now().Add(time.Hour), nil)
}

func (t *Task) run() {
	for {
		m := t.queue.DeleteMin()
		if m == nil {
			t.sleep()
			continue
		}
		d := m.Time.Sub(time.Now())
		select {
		case <-t.cha:
			t.queue.InsertAndSort(m)
		case <-time.After(d):
			if m.fun != nil {
				t.f.Puah(m.fun)
			}
		}
	}
}
