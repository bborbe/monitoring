package webdriver

import (
	monitoring_check "github.com/bborbe/monitoring/check"
)

type webdriverCheck struct{}

func New() *webdriverCheck {
	return new(webdriverCheck)
}

func (w *webdriverCheck) Check() monitoring_check.CheckResult {
	return nil
}

func (w *webdriverCheck) Description() string {
	return "foo"
}
