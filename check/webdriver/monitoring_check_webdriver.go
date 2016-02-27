package webdriver

import (
	"fmt"

	"time"

	"strings"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	"github.com/bborbe/webdriver"
)

const (
	DEFAULT_TIMEOUT = 30 * time.Second
)

var logger = log.DefaultLogger

type Action func(session *webdriver.Session) error

type webdriverCheck struct {
	url     string
	actions []Action
	timeout time.Duration
}

func New(url string) *webdriverCheck {
	w := new(webdriverCheck)
	w.url = url
	w.timeout = DEFAULT_TIMEOUT
	return w
}

func (w *webdriverCheck) Check() monitoring_check.CheckResult {
	logger.Debugf("webdriver check url %s", w.url)
	start := time.Now()
	return monitoring_check.NewCheckResult(w, w.check(), time.Now().Sub(start))
}

func (w *webdriverCheck) Description() string {
	return fmt.Sprintf("webdriver check on url %s", w.url)
}

func (w *webdriverCheck) Timeout(timeout time.Duration) *webdriverCheck {
	w.timeout = timeout
	return w
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

	if err = session.SetTimeouts("script", int(w.timeout)); err != nil {
		return err
	}

	logger.Debugf("fetch url")
	if err = session.Url(w.url); err != nil {
		return err
	}
	for _, action := range w.actions {
		if err = action(session); err != nil {
			return err
		}
	}
	return nil
}

func (h *webdriverCheck) AddAction(action Action) *webdriverCheck {
	h.actions = append(h.actions, action)
	return h
}

func (h *webdriverCheck) ExpectTitle(expectedTitle string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("expect title '%s' - started", expectedTitle)
		title, err := session.Title()
		if err != nil {
			logger.Debugf("expect title '%s' - failed", expectedTitle)
			return err
		}
		if !strings.Contains(title, expectedTitle) {
			logger.Debugf("expect title '%s' - failed", expectedTitle)
			return fmt.Errorf("expected title '%s' but got '%s'", expectedTitle, title)
		}
		logger.Debugf("expect title '%s' - success", expectedTitle)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Fill(xpath string, value string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("fill value '%s' to '%s' - started", value, xpath)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = session.FindElements(webdriver.XPath, xpath); err != nil {
			logger.Debugf("fill value '%s' to '%s' - failed", value, xpath)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("fill value '%s' to '%s' - failed", value, xpath)
			return fmt.Errorf("element '%s' not found", xpath)
		}
		for _, webElement := range webElements {
			if err = webElement.SendKeys(value); err != nil {
				logger.Debugf("fill value '%s' to '%s' - failed", value, xpath)
				return err
			}
		}
		logger.Debugf("fill value '%s' to '%s' - success", value, xpath)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Submit(xpath string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("submit '%s' - started", xpath)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = session.FindElements(webdriver.XPath, xpath); err != nil {
			logger.Debugf("submit '%s' - failed", xpath)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("submit '%s' - failed", xpath)
			return fmt.Errorf("element '%s' not found", xpath)
		}
		for _, webElement := range webElements {
			if err = webElement.Submit(); err != nil {
				logger.Debugf("submit '%s' - failed", xpath)
				return err
			}
		}
		logger.Debugf("submit '%s' - success", xpath)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Click(xpath string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("click '%s' - started", xpath)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = session.FindElements(webdriver.XPath, xpath); err != nil {
			logger.Debugf("click '%s' - failed", xpath)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("click '%s' - failed", xpath)
			return fmt.Errorf("element '%s' not found", xpath)
		}
		for _, webElement := range webElements {
			if err = webElement.Click(); err != nil {
				logger.Debugf("click '%s' - failed", xpath)
				return err
			}
		}
		logger.Debugf("click '%s' - success", xpath)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Exists(xpath string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("exists '%s' - started", xpath)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = session.FindElements(webdriver.XPath, xpath); err != nil {
			logger.Debugf("exists '%s' - failed", xpath)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("exists '%s' - failed", xpath)
			return fmt.Errorf("element '%s' not found", xpath)
		}
		logger.Debugf("exists '%s' - success", xpath)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) NotExists(xpath string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("notexists '%s' - started", xpath)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = session.FindElements(webdriver.XPath, xpath); err != nil {
			logger.Debugf("notexists '%s' - failed", xpath)
			return err
		}
		if len(webElements) != 0 {
			logger.Debugf("notexists '%s' - failed", xpath)
			return fmt.Errorf("element '%s' found", xpath)
		}
		logger.Debugf("notexists '%s' - success", xpath)
		return nil
	}
	h.AddAction(action)
	return h
}
