package task

import (
	"strings"
	"time"
)

var unix0 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)

// 产生固定间隔的时间定时函数
//  start:    开始时间
//  interval: 执行间隔
//  count:    执行次数 如果 -1 则不限制次数
func PeriodicIntervalCount(start time.Time, interval time.Duration, count int) func() time.Time {
	if start.IsZero() { // 开始时间 未初始化 则设置为 标准零点
		start = time.Unix(0, 0)
	}
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
func PeriodicEveryDay(tim string) func() time.Time {
	errRet := func() time.Time { return time.Time{} }
	sp := strings.SplitN(tim, " ", 2)
	switch len(sp) {
	case 1:
		tim = "1970-01-01 " + tim
	case 2:
	default:
		return errRet
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999", tim, time.Local)
	if err != nil {
		return errRet
	}
	return PeriodicIntervalCount(t, time.Hour*24, -1)
}
