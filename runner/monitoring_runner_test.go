package runner

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check"
)

type checkDummy struct {
	result      check.CheckResult
	description string
}

func TestRun(t *testing.T) {
	var err error
	checks := make([]check.Check, 0)
	checks = append(checks, NewCheckDummy(check.NewCheckResultSuccess("ok"), "ok"))
	results := Run(checks)
	for i := 0; i < len(checks); i++ {
		result := <-results
		err = AssertThat(result.Success(), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func NewCheckDummy(result check.CheckResult, description string) check.Check {
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
