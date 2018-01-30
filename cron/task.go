package cron

import (
	cr "github.com/robfig/cron"
)

var c *cr.Cron

func getCron() *cr.Cron {
	if c == nil {
		c = cr.New()
	}

	return c
}

var cIsStarted bool

func startCron() {
	if c == nil {
		return
	}

	if cIsStarted {
		return
	}

	cIsStarted = true
	c.Start()
}

/**
period: "0 30 * * * *", "@hourly", "@every 1h30m"
A cron expression represents a set of times, using 6 space-separated fields.
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/
func AddTask(period string, f func()) {
	c := getCron()

	c.AddFunc(period, f)

	startCron()
}
