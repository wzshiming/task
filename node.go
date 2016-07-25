package task

import (
	"llrb"
	"time"
)

type node struct {
	Time time.Time
	fun  func()
}

func (no *node) Less(i llrb.Item) bool {
	b, _ := i.(*node)
	if b == nil {
		return true
	}
	return no.Time.Before(b.Time)
}
