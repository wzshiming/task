package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/fork"
)

var TaskExit = time.Time{} // 间隔任务退出 标识

var unix0 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local) // 标准零点

var none = struct{}{} // 信号

type Task struct {
	fork  *fork.Fork    // 线程控制
	queue *list         // 任务队列
	curr  *Node         // 当前等待的
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
func (t *Task) Cancel(n *Node) {
	t.add(&Node{
		time: time.Unix(0, 0),
		task: func() {
			t.queue.Delete(n)
		},
	})
}

// 取消全部任务
func (t *Task) CancelAll(n *Node) {
	for _, v := range t.List() {
		t.Cancel(v)
	}
}

// 任务加入队列
func (t *Task) add(n *Node) *Node {
	select { // 判断管理线程是否运行 如果没有则启动
	case t.iru <- none:
		t.fork.Push(func() {
			t.run()
			<-t.iru
		})
	default:
	}

	t.queue.InsertAndSort(n) // 队列里插入

	if t.queue.Min() == n { // 如果插入到了第一个则刷新时间
		t.flash()
	}
	return n
}

// 新的任务
func (t *Task) Add(tim time.Time, task func()) *Node {
	return t.add(&Node{
		time: tim,
		task: func() {
			t.fork.Push(task)
		},
		name: ffmt.FMakeStack(1),
	})
}

// 重复任务加入队列
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

// 新的重复任务
func (t *Task) AddPeriodic(perfunc func() time.Time, task func()) (n *Node) {
	n = &Node{
		task: func() {
			t.addPeriodic(perfunc, n)
			t.fork.Push(task)
		},
		name: ffmt.FMakeStack(1),
	}
	return t.addPeriodic(perfunc, n)
}

// 刷新第一个执行的任务
func (t *Task) flash() {
	select {
	case t.ins <- none:
	default:
	}
}

// 不刷新
func (t *Task) unflash() {
	select {
	case <-t.ins:
	default:
	}
}

// 任务执行循环
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
		timer.Reset(sub)                   // 重置定时器
		select {
		case <-t.ins: // 有新的 任务节点插入
			t.queue.InsertAndSort(t.curr)
		case <-timer.C: // 到达最近执行的任务
			t.curr.task()
			t.unflash()
		}
	}
}

// 等待执行的任务数量
func (t *Task) Len() int {
	b := 0
	if t.curr != nil {
		b++
	}
	return t.queue.Len() + b
}

// 获取全部列表
func (t *Task) List() []*Node {
	if t.curr != nil {
		return append([]*Node{t.curr}, t.queue.List()...)
	}
	return t.queue.List()
}

// 打印出全部列表
func (t *Task) Print() error {
	sss := [][]string{{"NAME", "NEXT"}}
	for _, v := range t.List() {
		sss = append(sss, []string{v.Name(), v.Next().String()})
	}
	ss := ffmt.FmtTable(sss)
	_, err := fmt.Print(strings.Join(ss, "\n"), "\n")
	return err
}
