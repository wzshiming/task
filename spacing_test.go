package task

import (
	"testing"
	"time"
)

func TestSpacing(t *testing.T) {
	ta := NewTask(1)
	sp := NewSpacing(time.Second, func() {
		t.Log(time.Now())
	})

	for i := 0; i != 100; i++ {
		ta.Add(time.Now().Add(time.Second/10*time.Duration(i)), func() {
			sp.On()
		})
	}
	ta.Join()
}
