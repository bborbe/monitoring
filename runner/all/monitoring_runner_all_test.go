package all

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/check/dummy"
	"github.com/bborbe/monitoring/runner"
)

func TestImplementsRunner(t *testing.T) {
	c := New()
	var i *runner.Runner
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	var err error
	checks := make([]check.Check, 0)
	checks = append(checks, dummy.New(check.NewCheckResultSuccess("ok"), "ok"))
	results := Run(checks)
	for i := 0; i < len(checks); i++ {
		result := <-results
		err = AssertThat(result.Success(), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
}
