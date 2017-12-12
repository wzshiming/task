package task

import (
	"sort"
	"strings"
	"time"
)

// 产生固定间隔的时间定时函数
//  start:    开始时间
//  interval: 执行间隔
//  count:    执行次数 如果 -1 则不限制次数
func PeriodicIntervalCount(start time.Time, interval time.Duration, count int) func() time.Time {
	// 开始时间 未初始化 则设置为 标准零点
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

// 产生固定间隔的时间定时函数
//  offset:   执行时间的偏移
//  interval: 执行间隔
func PeriodicInterval(offset time.Duration, interval time.Duration) func() time.Time {
	return PeriodicIntervalCount(unix0.Add(offset), interval, -1)
}

// 每天的 某个时间执行
//  tim:
//   15:04:05 或 15:04:05.999999999 每天的这个时候
//   2006-01-02 15:04:05.999999999 从某个日期起 每天的这个时候
func ParseTimeDay(tim string) time.Time {
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

// 每天的 某个时间执行
//  tim:
//   15:04:05 或 15:04:05.999999999 每天的这个时候
//   2006-01-02 15:04:05.999999999 从某个日期起 每天的这个时候
func PeriodicEveryDay(tim string) func() time.Time {
	t := ParseTimeDay(tim)
	if t == TaskExit {
		return nil
	}
	return PeriodicIntervalCount(t, time.Hour*24, -1)
}

// 按近到远排序时间
type TimeSlice []time.Time

func (p TimeSlice) Len() int           { return len(p) }
func (p TimeSlice) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p TimeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TimeSlice) Sort()              { sort.Sort(p) }

// 产生固定间隔的时间定时函数
//  ts: 指定多个执行的时间
func PeriodicTiming(ts ...time.Time) func() time.Time {
	now := time.Now()
	// 排序执行的时间
	TimeSlice(ts).Sort()
	// 移除已经超过的时间
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
