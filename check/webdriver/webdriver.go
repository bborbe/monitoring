package webdriver

import (
	"fmt"

	"time"

	"strings"

	monitoring_check "github.com/bborbe/monitoring/check"
	"github.com/bborbe/webdriver"
	"github.com/golang/glog"
)

const (
	DEFAULT_TIMEOUT = 30 * time.Second
)

type Action func(session *webdriver.Session) error

type check struct {
	url       string
	actions   []Action
	webDriver webdriver.WebDriver
	timeout   time.Duration
}

func New(webDriver webdriver.WebDriver, url string) *check {
	c := new(check)
	c.url = url
	c.timeout = DEFAULT_TIMEOUT
	c.webDriver = webDriver
	return c
}

func (c *check) Check() monitoring_check.CheckResult {
	glog.V(2).Infof("webdriver check url %s", c.url)
	start := time.Now()
	return monitoring_check.NewCheckResult(c, c.check(), time.Now().Sub(start))
}

func (c *check) Description() string {
	return fmt.Sprintf("webdriver check on url %s", c.url)
}

func (c *check) Timeout(timeout time.Duration) *check {
	c.timeout = timeout
	return c
}

func (c *check) check() error {
	desired := webdriver.Capabilities{
		"Platform": "Linux",
		"phantomjs.page.customHeaders.Accept-Language": "en-US",
	}
	required := webdriver.Capabilities{}
	glog.V(2).Infof("create new session")
	session, err := c.webDriver.NewSession(desired, required)
	if err != nil {
		return err
	}
	defer session.Delete()

	if err = session.SetTimeouts("script", int(c.timeout)); err != nil {
		return err
	}

	glog.V(2).Infof("fetch url")
	if err = session.Url(c.url); err != nil {
		return err
	}
	for _, action := range c.actions {
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

func (c *check) AddAction(action Action) *check {
	c.actions = append(c.actions, action)
	return c
}

func (c *check) ExpectTitle(expectedTitle string) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("expect title '%s' - started", expectedTitle)
		title, err := session.Title()
		if err != nil {
			glog.V(2).Infof("expect title '%s' - failed", expectedTitle)
			return err
		}
		if !strings.Contains(title, expectedTitle) {
			glog.V(2).Infof("expect title '%s' - failed", expectedTitle)
			return fmt.Errorf("expected title '%s' but got '%s'", expectedTitle, title)
		}
		glog.V(2).Infof("expect title '%s' - success", expectedTitle)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) ExecuteScript(javascript string) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("execute script '%s' - started", javascript)
		args := make([]interface{}, 0)
		result, err := session.ExecuteScript(javascript, args)
		if err != nil {
			return err
		}
		glog.V(2).Infof("script result: %s", string(result))

		glog.V(2).Infof("execute script '%s' - success", javascript)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) Fill(strategy webdriver.FindElementStrategy, query string, value string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("fill value '%s' to '%s' - started", value, query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			glog.V(2).Infof("fill value '%s' to '%s' - failed", value, query)
			return err
		}
		if len(webElements) == 0 {
			glog.V(2).Infof("fill value '%s' to '%s' - failed", value, query)
			return fmt.Errorf("element '%s' not found", query)
		}
		for _, webElement := range webElements {
			if err = webElement.SendKeys(value); err != nil {
				glog.V(2).Infof("fill value '%s' to '%s' - failed", value, query)
				return err
			}
		}
		glog.V(2).Infof("fill value '%s' to '%s' - success", value, query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) Submit(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("submit '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			glog.V(2).Infof("submit '%s' - failed", query)
			return err
		}
		if len(webElements) == 0 {
			glog.V(2).Infof("submit '%s' - failed", query)
			return fmt.Errorf("element '%s' not found", query)
		}
		for _, webElement := range webElements {
			if err = webElement.Submit(); err != nil {
				glog.V(2).Infof("submit '%s' - failed", query)
				return err
			}
		}
		glog.V(2).Infof("submit '%s' - success", query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) Click(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("click '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			glog.V(2).Infof("click '%s' - failed", query)
			return err
		}
		if len(webElements) == 0 {
			glog.V(2).Infof("click '%s' - failed", query)
			return fmt.Errorf("element '%s' not found", query)
		}
		for _, webElement := range webElements {
			if err = webElement.Click(); err != nil {
				glog.V(2).Infof("click '%s' - failed", query)
				return err
			}
		}
		glog.V(2).Infof("click '%s' - success", query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) Exists(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("exists '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			glog.V(2).Infof("exists '%s' - failed", query)
			return err
		}
		if len(webElements) == 0 {
			glog.V(2).Infof("exists '%s' - failed", query)
			return fmt.Errorf("element '%s' not found", query)
		}
		glog.V(2).Infof("exists '%s' - success", query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) NotExists(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("notexists '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElementsNot(session, strategy, query, duration); err != nil {
			glog.V(2).Infof("notexists '%s' - failed", query)
			return err
		}
		if len(webElements) != 0 {
			glog.V(2).Infof("notexists '%s' - failed", query)
			return fmt.Errorf("element '%s' found", query)
		}
		glog.V(2).Infof("notexists '%s' - success", query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) PrintSource() *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("printsource - started")
		source, err := session.Source()
		if err != nil {
			glog.V(2).Infof("printsource  - failed")
			return fmt.Errorf("element found")
		}
		fmt.Println(source)
		glog.V(2).Infof("printsource - success")
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) WaitForDisplayed(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("wait for displayed '%s' - started", query)

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
			glog.V(2).Infof("wait for displayed '%s' - failed", query)
			return err
		}

		glog.V(2).Infof("wait for displayed '%s' - success", query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) WaitFor(strategy webdriver.FindElementStrategy, query string, duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("waitfor '%s' - started", query)
		var err error
		var webElements []webdriver.WebElement
		if webElements, err = findElements(session, strategy, query, duration); err != nil {
			return err
		}
		if len(webElements) == 0 {
			return fmt.Errorf("wait for element '%s' failed", query)
		}
		glog.V(2).Infof("waitfor '%s' - success", query)
		return nil
	}
	c.AddAction(action)
	return c
}

func (c *check) Sleep(duration time.Duration) *check {
	var action Action
	action = func(session *webdriver.Session) error {
		glog.V(2).Infof("sleep %dms - started", duration/time.Millisecond)
		time.Sleep(duration)
		glog.V(2).Infof("sleep %dms - success", duration/time.Millisecond)
		return nil
	}
	c.AddAction(action)
	return c
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
	glog.V(2).Infof("find elements")
	var err error
	start := time.Now()
	var exitConstraintError error
	for {
		var webElements []webdriver.WebElement
		glog.V(2).Infof("execute action")
		if webElements, err = action(); err != nil {
			glog.V(2).Infof("execute action failed")
			return nil, err
		}
		glog.V(2).Infof("check exit constraint")
		if exitConstraintError = exitConstraint(webElements); exitConstraintError == nil {
			glog.V(2).Infof("exit constraint succeed after %dms", time.Now().Sub(start)/time.Millisecond)
			return webElements, nil
		}
		if start.Add(duration).Before(time.Now()) {
			return nil, fmt.Errorf("exit constraint not succeed after %dms. %v", duration/time.Millisecond, exitConstraintError)
		}
		glog.V(2).Infof("exit constraint failed => sleep")
		time.Sleep(100 * time.Millisecond)
	}
}
