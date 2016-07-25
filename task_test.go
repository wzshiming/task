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

	time.Sleep(time.Second * 2)
}

func TestB(t *testing.T) {
	task := NewTask(10)
	fmt.Println("000", time.Now())
	task.AddPeriodic(PeriodicIntervalCount(time.Now().Add(time.Second*1), time.Second, 3), func() {
		fmt.Println("111", time.Now())
	})

	time.Sleep(time.Second * 7)

}
