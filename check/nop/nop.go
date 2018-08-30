package nop

import (
	"time"

	monitoring_check "github.com/bborbe/monitoring/check"
)

type check struct {
	name string
}

func New(name string) monitoring_check.Check {
	c := new(check)
	c.name = name
	return c
}

func (c *check) Check() monitoring_check.CheckResult {
	start := time.Now()
	return monitoring_check.NewCheckResultSuccess(c.name, time.Now().Sub(start))
}

func (c *check) Description() string {
	return c.name
}
