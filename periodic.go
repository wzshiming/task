package task

import "time"

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
	return PeriodicIntervalCount(time.Unix(0, 0).Add(offset), interval, -1)
}
