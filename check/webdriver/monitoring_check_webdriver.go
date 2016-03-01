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
	url       string
	actions   []Action
	webDriver webdriver.WebDriver
	timeout   time.Duration
}

func New(webDriver webdriver.WebDriver, url string) *webdriverCheck {
	w := new(webdriverCheck)
	w.url = url
	w.timeout = DEFAULT_TIMEOUT
	w.webDriver = webDriver
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
	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	logger.Debugf("create new session")
	session, err := w.webDriver.NewSession(desired, required)
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
			//printSource(session)
			return err
		}
	}
	return nil
}

func printSource(session *webdriver.Session) {
	source, _ := session.Source()
	fmt.Printf("source:\n%s\n", source)
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

func (h *webdriverCheck) ExecuteScript(javascript string) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("execute script '%s' - started", javascript)
		args := make([]interface{}, 0)
		result, err := session.ExecuteScript(javascript, args)
		if err != nil {
			return err
		}
		logger.Debugf("script result: %s", string(result))

		logger.Debugf("execute script '%s' - success", javascript)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Fill(strategy webdriver.FindElementStrategy, query string, value string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("fill value '%s' to '%s' - started", value, query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			logger.Debugf("fill value '%s' to '%s' - failed", value, query)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("fill value '%s' to '%s' - failed", value, query)
			return fmt.Errorf("element '%s' not found", query)
		}
		for _, webElement := range webElements {
			if err = webElement.SendKeys(value); err != nil {
				logger.Debugf("fill value '%s' to '%s' - failed", value, query)
				return err
			}
		}
		logger.Debugf("fill value '%s' to '%s' - success", value, query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Submit(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("submit '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			logger.Debugf("submit '%s' - failed", query)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("submit '%s' - failed", query)
			return fmt.Errorf("element '%s' not found", query)
		}
		for _, webElement := range webElements {
			if err = webElement.Submit(); err != nil {
				logger.Debugf("submit '%s' - failed", query)
				return err
			}
		}
		logger.Debugf("submit '%s' - success", query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Click(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("click '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			logger.Debugf("click '%s' - failed", query)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("click '%s' - failed", query)
			return fmt.Errorf("element '%s' not found", query)
		}
		for _, webElement := range webElements {
			if err = webElement.Click(); err != nil {
				logger.Debugf("click '%s' - failed", query)
				return err
			}
		}
		logger.Debugf("click '%s' - success", query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Exists(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("exists '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			logger.Debugf("exists '%s' - failed", query)
			return err
		}
		if len(webElements) == 0 {
			logger.Debugf("exists '%s' - failed", query)
			return fmt.Errorf("element '%s' not found", query)
		}
		logger.Debugf("exists '%s' - success", query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) NotExists(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("notexists '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElementsNot(session, strategy, query, duration); err != nil {
			logger.Debugf("notexists '%s' - failed", query)
			return err
		}
		if len(webElements) != 0 {
			logger.Debugf("notexists '%s' - failed", query)
			return fmt.Errorf("element '%s' found", query)
		}
		logger.Debugf("notexists '%s' - success", query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) PrintSource() *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("printsource - started")
		source, err := session.Source()
		if err != nil {
			logger.Debugf("printsource  - failed")
			return fmt.Errorf("element found")
		}
		fmt.Println(source)
		logger.Debugf("printsource - success")
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) WaitForDisplayed(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("wait for displayed '%s' - started", query)

		_, err := findElementsWait(func() ([]webdriver.WebElement, error) {
			return session.FindElements(strategy, query)
		}, func(webElements []webdriver.WebElement) error {
			for _, webElement := range webElements {
				displayed, err := webElement.IsDisplayed()
				if err != nil {
					return err
				}
				if !displayed {
					return fmt.Errorf("element '%s' found but not displayed", query)
				}
			}
			if len(webElements) == 0 {
				return fmt.Errorf("element '%s' not found", query)
			}
			return nil
		}, duration)

		if err != nil {
			logger.Debugf("wait for displayed '%s' - failed", query)
			return err
		}

		logger.Debugf("wait for displayed '%s' - success", query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) WaitFor(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("waitfor '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			return err
		}
		if len(webElements) == 0 {
			return fmt.Errorf("wait for element '%s' failed", query)
		}
		logger.Debugf("waitfor '%s' - success", query)
		return nil
	}
	h.AddAction(action)
	return h
}

func (h *webdriverCheck) Sleep(duration time.Duration) *webdriverCheck {
	var action Action
	action = func(session *webdriver.Session) error {
		logger.Debugf("sleep %dms - started", duration/time.Millisecond)
		time.Sleep(duration)
		logger.Debugf("sleep %dms - success", duration/time.Millisecond)
		return nil
	}
	h.AddAction(action)
	return h
}

func findElements(session *webdriver.Session, strategy webdriver.FindElementStrategy, query string, duration time.Duration) ([]webdriver.WebElement, error) {
	return findElementsWait(func() ([]webdriver.WebElement, error) {
		return session.FindElements(strategy, query)
	}, func(webElements []webdriver.WebElement) error {
		if len(webElements) == 0 {
			return fmt.Errorf("element '%s' not found", query)
		}
		return nil
	}, duration)
}

func findElementsNot(session *webdriver.Session, strategy webdriver.FindElementStrategy, query string, duration time.Duration) ([]webdriver.WebElement, error) {
	return findElementsWait(func() ([]webdriver.WebElement, error) {
		return session.FindElements(strategy, query)
	}, func(webElements []webdriver.WebElement) error {
		if len(webElements) != 0 {
			return fmt.Errorf("element '%s' found", query)
		}
		return nil
	}, duration)
}

func findElementsWait(action func() ([]webdriver.WebElement, error), exitConstraint func([]webdriver.WebElement) error, duration time.Duration) ([]webdriver.WebElement, error) {
	logger.Debugf("find elements")
	var err error
	start := time.Now()
	var exitConstraintError error
	for {
		var webElements []webdriver.WebElement
		logger.Debugf("execute action")
		if webElements, err = action(); err != nil {
			logger.Debugf("execute action failed")
			return nil, err
		}
		logger.Debugf("check exit constraint")
		if exitConstraintError = exitConstraint(webElements); exitConstraintError == nil {
			logger.Debugf("exit constraint succeed after %dms", time.Now().Sub(start)/time.Millisecond)
			return webElements, nil
		}
		if start.Add(duration).Before(time.Now()) {
			return nil, fmt.Errorf("exit constraint not succeed after %dms. %v", duration/time.Millisecond, exitConstraintError)
		}
		logger.Debugf("exit constraint failed => sleep")
		time.Sleep(100 * time.Millisecond)
	}
}
