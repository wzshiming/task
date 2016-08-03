package task

import (
	"time"

	"github.com/wzshiming/fork"
)

type Task struct {
	queue *List
	cha   chan struct{}
	f     *fork.Fork
}

// 任务管理
//  i: 线程数
func NewTask(i int) *Task {
	t := &Task{
		queue: NewList(),
		cha:   make(chan struct{}, 1),
		f:     fork.NewFork(i),
	}
	go t.run()
	return t
}

// 结束任务
func (t *Task) Close(n *node) {
	// 把可能在执行第一个等待中的 待删除 挤出去
	t.Add(time.Now().Add(-time.Second), func() {
		t.queue.Delete(n)
	})
	// 防止 被挤出去 马上又回来了
	t.Add(time.Now().Add(-time.Second/2), nil)
}

// 任务加入队列
func (t *Task) add(n *node) *node {
	t.queue.InsertAndSort(n)
	if t.queue.Min() == n {
		t.flash()
	}
	return n
}

// 新的任务
func (t *Task) Add(tim time.Time, f func()) *node {
	return t.add(&node{
		tim: tim,
		fun: f,
	})
}

// 重复任务加入对了
func (t *Task) addPeriodic(perfunc func() time.Time, n *node) *node {
	p := perfunc()
	if p.IsZero() {
		return nil
	}
	n.tim = p
	return t.add(n)
}

// 新的重复任务
func (t *Task) AddPeriodic(perfunc func() time.Time, f func()) (n *node) {
	n = &node{
		fun: func() {
			t.addPeriodic(perfunc, n)
			f()
		},
	}
	return t.addPeriodic(perfunc, n)
}

// 刷新第一个执行的任务
func (t *Task) flash() {
	select {
	case t.cha <- struct{}{}:
	default:
	}
}

// 没有任务时休眠
func (t *Task) sleep() {
	t.Add(time.Now().Add(time.Hour), nil)
}

// 执行单次任务
func (t *Task) exec(n *node) {
	select {
	case <-t.cha:
		if n.fun != nil {
			t.queue.InsertAndSort(n)
		}
	case <-time.After(n.tim.Sub(time.Now())):
		if n.fun != nil {
			t.f.Puah(n.fun)
		}
	}
}

// 任务执行循环
func (t *Task) run() {
	for {
		n := t.queue.DeleteMin()
		if n == nil {
			t.sleep()
			continue
		}
		t.exec(n)
	}
}

// 等待执行的任务数量 不算第一个
func (t *Task) Len() int {
	return t.queue.Len()
}
