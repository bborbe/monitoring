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
	url             string
	expectedTitle   string
	expectedContent string
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
	err = h.checkContent(content)
	if err != nil {
		return check.NewCheckResult(h, err)
	}
	return check.NewCheckResult(h, err)
}

func (h *httpCheck) ExpectTitle(expectedTitle string) *httpCheck {
	h.expectedTitle = expectedTitle
	return h
}
func (h *httpCheck) ExpectContent(expectedContent string) *httpCheck {
	h.expectedContent = expectedContent
	return h
}

func (h *httpCheck) checkContent(content []byte) error {
	return checkContent(h.expectedContent, content)
}

func checkContent(expectedContent string, content []byte) error {
	if len(expectedContent) == 0 {
		return nil
	}
	logger.Debugf("content: %s", string(content))
	expression := fmt.Sprintf(`(?is)<html[^>]*>.*?<body[^>]*>.*?%s.*?</body>.*?</html>`, regexp.QuoteMeta(expectedContent))
	logger.Debugf("regexp: %s", expression)
	re := regexp.MustCompile(expression)
	if len(re.FindSubmatch(content)) > 0 {
		return nil
	}
	return fmt.Errorf("content %s not found", expectedContent)
}

func (h *httpCheck) checkTitle(content []byte) error {
	return checkTitle(h.expectedTitle, content)
}

func checkTitle(expectedTitle string, content []byte) error {
	if len(expectedTitle) == 0 {
		return nil
	}
	logger.Debugf("content: %s", string(content))
	expression := fmt.Sprintf(`(?is)<html[^>]*>.*?<head[^>]*>.*?<title[^>]*>%s</title>.*?</head>.*?</html>`, regexp.QuoteMeta(expectedTitle))
	logger.Debugf("regexp: %s", expression)
	re := regexp.MustCompile(expression)
	if len(re.FindSubmatch(content)) > 0 {
		return nil
	}
	return fmt.Errorf("title %s not found", expectedTitle)
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
