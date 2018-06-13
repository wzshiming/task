package task

import (
	"sort"
	"strings"
	"time"
)

// PeriodicIntervalCount is generates a fixed interval time function
func PeriodicIntervalCount(start time.Time, interval time.Duration, count int) func() time.Time {
	// If the start time is not initialized, it is set to standard zero
	if start.IsZero() {
		start = time.Unix(0, 0)
	}
	return func() time.Time {
		now := time.Now()
		sub := now.Sub(start)
		if count >= 0 && int(sub/interval) >= count {
			return TaskExit
		}
		if start.After(now) {
			return start
		}
		return now.Add(interval - sub%interval)
	}
}

// PeriodicInterval is generates a fixed interval time function, unlimited number of times
func PeriodicInterval(offset time.Duration, interval time.Duration) func() time.Time {
	return PeriodicIntervalCount(unix0.Add(offset), interval, -1)
}

// parseTimeDay is parse fixed time
func parseTimeDay(tim string) time.Time {
	sp := strings.SplitN(tim, " ", 2)
	switch len(sp) {
	case 1:
		tim = "1970-01-01 " + tim
	case 2:
	default:
		return TaskExit
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999", tim, time.Local)
	if err != nil {
		return TaskExit
	}
	return t
}

// PeriodicEveryDay is a fixed time of day
//   15:04:05 of 15:04:05.999999999 This time of day
//   2006-01-02 15:04:05.999999999 This time of day from a certain date
func PeriodicEveryDay(tim string) func() time.Time {
	t := parseTimeDay(tim)
	if t == TaskExit {
		return nil
	}
	return PeriodicIntervalCount(t, time.Hour*24, -1)
}

// TimeSlice sorted from near to far
type TimeSlice []time.Time

func (p TimeSlice) Len() int           { return len(p) }
func (p TimeSlice) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p TimeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TimeSlice) Sort()              { sort.Sort(p) }

// PeriodicTiming is multiple time for execution
func PeriodicTiming(ts ...time.Time) func() time.Time {
	now := time.Now()
	// sorted from near to far
	TimeSlice(ts).Sort()
	// remove time that has expired
	for _, v := range ts {
		if !v.Before(now) {
			break
		}
		ts = ts[1:]
	}
	//
	return func() time.Time {
		if len(ts) == 0 {
			return TaskExit
		}
		t := ts[0]
		ts = ts[1:]
		return t
	}
}
