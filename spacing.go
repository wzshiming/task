package task

import (
	"time"
)

// Spacing is interval to perform
type Spacing struct {
	ms      chan int
	ct      <-chan time.Time
	perfunc func() time.Time
	fun     func()
}

// NewSpacing create a fixed interval
func NewSpacing(d time.Duration, fun func()) *Spacing {
	return NewSpacingPeriodic(PeriodicIntervalCount(time.Now(), d, -1), fun)
}

// NewSpacingPeriodic create a compute interval
func NewSpacingPeriodic(perfunc func() time.Time, fun func()) *Spacing {
	s := &Spacing{
		ms:      make(chan int, 1),
		fun:     fun,
		perfunc: perfunc,
	}
	s.ct = s.getAfter()
	return s
}

func (s *Spacing) getAfter() <-chan time.Time {
	p := s.perfunc()
	now := time.Now()
	x := p.Sub(now)
	if x < 0 {
		x = 0
	}
	return time.After(x)
}

// On is trying to perform
func (s *Spacing) On() {
	select {
	case <-s.ct:
		s.fun()
		s.ct = s.getAfter()
	default:
	}
}
