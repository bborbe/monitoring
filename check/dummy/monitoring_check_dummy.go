package dummy

import (
	monitoring_check "github.com/bborbe/monitoring/check"
)

type checkDummy struct {
	result      monitoring_check.CheckResult
	description string
}

func New(result monitoring_check.CheckResult, description string) monitoring_check.Check {
	c := new(checkDummy)
	c.result = result
	c.description = description
	return c
}

func (c *checkDummy) Check() monitoring_check.CheckResult {
	return c.result
}

func (c *checkDummy) Description() string {
	return c.description
}
