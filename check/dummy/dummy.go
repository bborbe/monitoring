package dummy

import (
	monitoring_check "github.com/bborbe/monitoring/check"
)

type check struct {
	result      monitoring_check.CheckResult
	description string
}

func New(result monitoring_check.CheckResult, description string) monitoring_check.Check {
	c := new(check)
	c.result = result
	c.description = description
	return c
}

func (c *check) Check() monitoring_check.CheckResult {
	return c.result
}

func (c *check) Description() string {
	return c.description
}
