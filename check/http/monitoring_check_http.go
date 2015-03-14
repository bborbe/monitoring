package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"regexp"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
)

type ContentExpectation func([]byte) error

type httpCheck struct {
	url                 string
	contentExpectations []ContentExpectation
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
		logger.Debugf("fetch url failed %s: %v", h.url, err)
		return check.NewCheckResult(h, err)
	}
	for _, contentExpectation := range h.contentExpectations {
		err = contentExpectation(content)
		if err != nil {
			return check.NewCheckResult(h, err)
		}
	}
	return check.NewCheckResult(h, err)
}

func (h *httpCheck) AddExpectation(contentExpectation ContentExpectation) *httpCheck {
	h.contentExpectations = append(h.contentExpectations, contentExpectation)
	return h
}

func (h *httpCheck) ExpectTitle(expectedTitle string) *httpCheck {
	var contentExpectation ContentExpectation
	contentExpectation = func(content []byte) error {
		return checkTitle(expectedTitle, content)
	}
	h.AddExpectation(contentExpectation)
	return h
}
func (h *httpCheck) ExpectContent(expectedContent string) *httpCheck {
	var contentExpectation ContentExpectation
	contentExpectation = func(content []byte) error {
		return checkContent(expectedContent, content)
	}
	h.AddExpectation(contentExpectation)
	return h
}

func (h *httpCheck) ExpectBody(expectedBody string) *httpCheck {
	var contentExpectation ContentExpectation
	contentExpectation = func(content []byte) error {
		return checkBody(expectedBody, content)
	}
	h.AddExpectation(contentExpectation)
	return h
}

func checkContent(expectedContent string, content []byte) error {
	if len(expectedContent) == 0 {
		return nil
	}
	logger.Tracef("content: %s", string(content))
	expression := fmt.Sprintf(`(?is).*?%s.*?`, regexp.QuoteMeta(expectedContent))
	logger.Tracef("content regexp: %s", expression)
	re := regexp.MustCompile(expression)
	if len(re.FindSubmatch(content)) > 0 {
		return nil
	}
	return fmt.Errorf("content %s not found", expectedContent)
}

func checkBody(expectedBody string, content []byte) error {
	if len(expectedBody) == 0 {
		return nil
	}
	logger.Tracef("content: %s", string(content))
	expression := fmt.Sprintf(`(?is)<html[^>]*>.*?<body[^>]*>.*?%s.*?</body>.*?</html>`, regexp.QuoteMeta(expectedBody))
	logger.Tracef("body regexp: %s", expression)
	re := regexp.MustCompile(expression)
	if len(re.FindSubmatch(content)) > 0 {
		return nil
	}
	return fmt.Errorf("content %s not found", expectedBody)
}

func checkTitle(expectedTitle string, content []byte) error {
	if len(expectedTitle) == 0 {
		return nil
	}
	logger.Tracef("content: %s", string(content))
	expression := fmt.Sprintf(`(?is)<html[^>]*>.*?<head[^>]*>.*?<title[^>]*>\s*%s\s*</title>.*?</head>.*?</html>`, regexp.QuoteMeta(expectedTitle))
	logger.Tracef("title regexp: %s", expression)
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
