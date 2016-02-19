package all

import (
	"testing"

	"runtime"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_check_dummy "github.com/bborbe/monitoring/check/dummy"
	monitoring_runner "github.com/bborbe/monitoring/runner"
)

func TestImplementsRunner(t *testing.T) {
	c := New(123)
	var i *monitoring_runner.Runner
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	var err error
	checks := make([]monitoring_check.Check, 0)
	checks = append(checks, monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("ok"), "ok"))
	results := Run(runtime.NumCPU()*2, checks)
	for i := 0; i < len(checks); i++ {
		result := <-results
		err = AssertThat(result.Success(), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
}
