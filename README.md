# task

timed task

[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/task)](https://goreportcard.com/report/github.com/wzshiming/task)
[![GoDoc](https://godoc.org/github.com/wzshiming/task?status.svg)](https://godoc.org/github.com/wzshiming/task)
[![GitHub license](https://img.shields.io/github/license/wzshiming/task.svg)](https://github.com/wzshiming/task/blob/master/LICENSE)

- [English](https://github.com/wzshiming/task/blob/master/README.md)
- [简体中文](https://github.com/wzshiming/task/blob/master/README_cn.md)

## Usage

[API Documentation](https://godoc.org/github.com/wzshiming/task)

## example

``` go
package main

import (
	"time"

	"github.com/wzshiming/task"
)

func main() {
	// The maximum is 1 thread.
	t := task.NewTask(1)

	// Execute only once. Execute one second later.
	t.Add(time.Now().Add(time.Second), func() {})

	// Execute once per second.
	t.AddPeriodic(task.PeriodicInterval(0, time.Second), func() {})

	// Execute once per second. Use crontab definitions. Accurate to seconds.
	t.AddPeriodic(task.PeriodicCrontab("* * * * * *"), func() {})

	t.Join()
}

```

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/task/blob/master/LICENSE) for the full license text.
