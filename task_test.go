package task

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	task := NewTask(10)
	rand.Seed(time.Now().UnixNano())
	lll := 10000
	ccc := make(chan bool, lll)
	for i := 0; i != lll; i++ {
		//go func() {
		v := rand.Int()%1000 + 1
		d := time.Millisecond * time.Duration(v)

		task.Add(time.Now().Add(d), func() {
			fmt.Println("hello", d)
			ccc <- true
		})
		//}()
	}

	task.Join()
}

func TestB(t *testing.T) {
	task := NewTask(10)
	fmt.Println("begin", time.Now())
	task.AddPeriodic(PeriodicIntervalCount(time.Now(), time.Second/10, 50), func() {
		fmt.Println("fork", time.Now())
	})
	task.Join()
}

func TestC(t *testing.T) {
	task := NewTask(10)
	fmt.Println("begin", time.Now())
	n := task.AddPeriodic(PeriodicIntervalCount(time.Now(), time.Second, 5), func() {
		fmt.Println("fork", time.Now())
	})
	task.Add(time.Now().Add(time.Second*2), func() {
		task.Cancel(n)
		fmt.Println("close", time.Now())
	})
	task.Join()
	fmt.Println("begin", time.Now())
	n = task.AddPeriodic(PeriodicIntervalCount(time.Now(), time.Second, 5), func() {
		fmt.Println("fork", time.Now())
	})
	task.Add(time.Now().Add(time.Second*2), func() {
		task.Cancel(n)
		fmt.Println("close", time.Now())
	})
	task.Join()
}
