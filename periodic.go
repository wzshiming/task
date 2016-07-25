package task

import "time"

func PeriodicIntervalCount(start time.Time, interval time.Duration, count int) func() time.Time {
	return func() time.Time {
		now := time.Now()
		sub := now.Sub(start)
		if count >= 0 && sub/interval+1 >= time.Duration(count) {
			return time.Time{}
		}
		if start.After(now) {
			return start
		}
		return now.Add(interval - sub%interval)
	}
}
