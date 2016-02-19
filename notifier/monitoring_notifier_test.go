package notifier

import (
	"fmt"
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
)

func TestBuildMailContentNoResults(t *testing.T) {
	var err error
	results := make([]monitoring_check.CheckResult, 0)
	content := buildMailContent(results)
	err = AssertThat(content, Is("Checks executed: 0\nChecks failed: 0\n"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildMailContentSuccess(t *testing.T) {
	var err error
	results := make([]monitoring_check.CheckResult, 0)
	results = append(results, monitoring_check.NewCheckResultSuccess("ok"))
	content := buildMailContent(results)
	err = AssertThat(content, Is("Checks executed: 1\nChecks failed: 0\n"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildMailContentFail(t *testing.T) {
	var err error
	results := make([]monitoring_check.CheckResult, 0)
	results = append(results, monitoring_check.NewCheckResultFail("fail", fmt.Errorf("error")))
	content := buildMailContent(results)
	err = AssertThat(content, Is("Checks executed: 1\nChecks failed: 1\nfail - error\n"))
	if err != nil {
		t.Fatal(err)
	}
}
