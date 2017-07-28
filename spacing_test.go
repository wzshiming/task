package task

import (
	"testing"
	"time"

	"github.com/wzshiming/ffmt"
)

func TestSpacing(t *testing.T) {
	ta := NewTask(1)
	sp := NewSpacing(time.Second, func() {
		ffmt.Mark(time.Now())
	})

	for i := 0; i != 100; i++ {
		ta.Add(time.Now().Add(time.Second/10*time.Duration(i)), func() {
			sp.On()
		})
	}
	ta.Join()
}
