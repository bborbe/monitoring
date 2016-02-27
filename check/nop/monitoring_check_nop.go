package nop

import (
	"time"

	monitoring_check "github.com/bborbe/monitoring/check"
)

type checkNop struct {
	name string
}

func New(name string) monitoring_check.Check {
	c := new(checkNop)
	c.name = name
	return c
}

func (c *checkNop) Check() monitoring_check.CheckResult {
	start := time.Now()
	return monitoring_check.NewCheckResultSuccess("nop", time.Now().Sub(start))
}

func (c *checkNop) Description() string {
	return c.name
}
