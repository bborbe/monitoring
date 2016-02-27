package webdriver

import (
	"fmt"

	"time"

	"strings"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	"github.com/bborbe/webdriver"
)

var logger = log.DefaultLogger

type Expectation func(session *webdriver.Session) error

type webdriverCheck struct {
	url          string
	expectations []Expectation
}

func New(url string) *webdriverCheck {
	w := new(webdriverCheck)
	w.url = url
	return w
}

func (w *webdriverCheck) Check() monitoring_check.CheckResult {
	logger.Debugf("webdriver check url %s", w.url)
	start := time.Now()
	return monitoring_check.NewCheckResult(w, w.check(), time.Now().Sub(start))
}

func (h *webdriverCheck) Description() string {
	return fmt.Sprintf("webdriver check on url %s", h.url)
}

func (w *webdriverCheck) check() error {
	var err error
	logger.Debugf("create new driver")
	driver := webdriver.NewPhantomJsDriver("/opt/phantomjs-2.1.1-macosx/bin/phantomjs")
	if err = driver.Start(); err != nil {
		return err
	}
	defer driver.Stop()

	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	logger.Debugf("create new session")
	session, err := driver.NewSession(desired, required)
	if err != nil {
		return err
	}
	defer session.Delete()

	logger.Debugf("fetch url")
	if err = session.Url(w.url); err != nil {
		return err
	}

	for _, expectation := range w.expectations {
		if err = expectation(session); err != nil {
			return err
		}
	}
	return nil
}

func (h *webdriverCheck) AddExpectation(expectation Expectation) *webdriverCheck {
	h.expectations = append(h.expectations, expectation)
	return h
}

func (h *webdriverCheck) ExpectTitle(expectedTitle string) *webdriverCheck {
	var expectation Expectation
	expectation = func(session *webdriver.Session) error {
		title, err := session.Title()
		if err != nil {
			return err
		}
		if !strings.Contains(title, expectedTitle) {
			return fmt.Errorf("expected title '%s' but got '%s'", expectedTitle, title)
		}
		return nil
	}
	h.AddExpectation(expectation)
	return h
}
