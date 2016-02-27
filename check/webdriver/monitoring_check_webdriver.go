package webdriver

import (
	"fmt"

	monitoring_check "github.com/bborbe/monitoring/check"
)

type webdriverCheck struct {
	url string
}

func New(url string) *webdriverCheck {
	w := new(webdriverCheck)
	w.url = url
	return w
}

func (w *webdriverCheck) Check() monitoring_check.CheckResult {
	return nil
}

func (h *webdriverCheck) Description() string {
	return fmt.Sprintf("webdriver check on url %s", h.url)
}
