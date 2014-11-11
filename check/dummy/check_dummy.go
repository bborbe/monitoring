package dummy

import (
	"github.com/bborbe/monitoring/check"
)

type checkDummy struct {
	result      check.CheckResult
	description string
}

func New(result check.CheckResult, description string) check.Check {
	c := new(checkDummy)
	c.result = result
	c.description = description
	return c
}

func (c *checkDummy) Check() check.CheckResult {
	return c.result
}

func (c *checkDummy) Description() string {
	return c.description
}
