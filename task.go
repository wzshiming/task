package task

import (
	"time"

	"github.com/wzshiming/fork"
)

type Task struct {
	fork  *fork.Fork    // 线程控制
	queue *List         // 任务队列
	ins   chan struct{} // 插入新的任务的信号
	iru   chan struct{} // 管理线程是否运行中的信号
}

// 任务管理
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

// 直到 当前所有任务 完结
func (t *Task) Join() {
	t.fork.Join()
}

// 取消任务
func (t *Task) Cancel(n *node) {
	t.add(&node{
		time: time.Unix(0, 0),
		prefix: func() {
			t.queue.Delete(n)
		},
	})
}

// 任务加入队列
func (t *Task) add(n *node) *node {
	if len(t.iru) == 0 { // 判断管理线程是否运行 如果没有则启动
		t.fork.Puah(t.run)
	}
	t.queue.InsertAndSort(n)
	if t.queue.Min() == n {
		t.flash()
	}
	return n
}

// 新的任务
func (t *Task) Add(tim time.Time, f func()) *node {
	return t.add(&node{
		time: tim,
		task: f,
	})
}

// 重复任务加入队列
func (t *Task) addPeriodic(perfunc func() time.Time, n *node) *node {
	p := perfunc()
	if p.IsZero() {
		return nil
	}
	n.time = p
	return t.add(n)
}

// 新的重复任务
func (t *Task) AddPeriodic(perfunc func() time.Time, task func()) (n *node) {
	n = &node{
		prefix: func() {
			t.addPeriodic(perfunc, n)
		},
		task: task,
	}
	return t.addPeriodic(perfunc, n)
}

// 刷新第一个执行的任务
func (t *Task) flash() {
	select {
	case t.ins <- struct{}{}:
	default:
	}
}

// 任务执行循环
func (t *Task) run() {
	t.iru <- struct{}{}
	timer := time.NewTimer(time.Hour)
	for {
		n := t.queue.DeleteMin()
		if n == nil { // 如果没有任务了 结束线程
			if t.Len() == 0 {
				break
			}
			continue
		}
		sub := n.time.Sub(time.Now()) // 计算 休眠时长
		timer.Reset(sub)              // 重置定时器
		select {
		case <-t.ins: // 有新的 任务节点插入
			if n.task != nil {
				t.queue.InsertAndSort(n)
			}
		case <-timer.C: // 到达最近执行的任务
			if n.prefix != nil { // 前缀执行
				n.prefix()
			}
			if n.task != nil {
				t.fork.Puah(n.task)
			}
		}
	}
	<-t.iru
}

// 等待执行的任务数量 不算第一个
func (t *Task) Len() int {
	return t.queue.Len()
}
