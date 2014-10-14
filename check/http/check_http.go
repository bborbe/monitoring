package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"regexp"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
)

type httpCheck struct {
	url   string
	title string
}

var logger = log.DefaultLogger

func New(url string) *httpCheck {
	h := new(httpCheck)
	h.url = url
	return h
}

func (h *httpCheck) Description() string {
	return fmt.Sprintf("http check on url %s", h.url)
}

func (h *httpCheck) Check() check.CheckResult {
	content, err := get(h.url)
	if err != nil {
		return check.NewCheckResult(h, err)
	}
	err = h.checkTitle(content)
	if err != nil {
		return check.NewCheckResult(h, err)
	}
	return check.NewCheckResult(h, err)
}

func (h *httpCheck) checkTitle(content []byte) error {
	return checkTitle(h.title, content)
}

func (h *httpCheck) ExpectTitle(title string) *httpCheck {
	h.title = title
	return h
}

func checkTitle(title string, content []byte) error {
	if len(title) == 0 {
		return nil
	}
	logger.Debugf("content: %s", string(content))
	expression := fmt.Sprintf(`(?is)<html[^>]*>.*?<head[^>]*>.*?<title[^>]*>%s</title>.*?</head>.*?</html>`, regexp.QuoteMeta(title))
	logger.Debugf("regexp: %s", expression)
	re := regexp.MustCompile(expression)
	if len(re.FindSubmatch(content)) > 0 {
		return nil
	}

	return fmt.Errorf("title %s not found", title)
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
