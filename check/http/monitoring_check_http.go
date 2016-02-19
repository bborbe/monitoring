package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/redirect_follower"
	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type Expectation func(httpResponse *HttpResponse) error

type httpCheck struct {
	url            string
	username       string
	password       string
	passwordFile   string
	expectations   []Expectation
	executeRequest ExecuteRequest
}

type HttpResponse struct {
	Content    []byte
	StatusCode int
}

var logger = log.DefaultLogger

func New(url string) *httpCheck {
	h := new(httpCheck)
	h.url = url
	redirectFollower := redirect_follower.New(http_client_builder.New().WithoutProxy().BuildRoundTripper().RoundTrip)
	h.executeRequest = redirectFollower.ExecuteRequestAndFollow
	return h
}

func (h *httpCheck) Description() string {
	return fmt.Sprintf("http check on url %s", h.url)
}

func (h *httpCheck) Check() monitoring_check.CheckResult {
	if len(h.password) == 0 && len(h.passwordFile) > 0 {
		logger.Debugf("read password from file %s", h.passwordFile)
		password, err := ioutil.ReadFile(h.passwordFile)
		if err != nil {
			logger.Debugf("read password file failed %s: %v", h.passwordFile, err)
			return monitoring_check.NewCheckResult(h, err)
		}
		h.password = strings.TrimSpace(string(password))
	}
	httpResponse, err := get(h.executeRequest, h.url, h.username, h.password)
	if err != nil {
		logger.Debugf("fetch url failed %s: %v", h.url, err)
		return monitoring_check.NewCheckResult(h, err)
	}
	for _, expectation := range h.expectations {
		if err = expectation(httpResponse); err != nil {
			return monitoring_check.NewCheckResult(h, err)
		}
	}
	return monitoring_check.NewCheckResult(h, err)
}

func (h *httpCheck) AddExpectation(expectation Expectation) *httpCheck {
	h.expectations = append(h.expectations, expectation)
	return h
}

func (h *httpCheck) ExpectTitle(expectedTitle string) *httpCheck {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkTitle(expectedTitle, resp.Content)
	}
	h.AddExpectation(expectation)
	return h
}

func (h *httpCheck) ExpectStatusCode(expectedStatusCode int) *httpCheck {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkStatusCode(expectedStatusCode, resp.StatusCode)
	}
	h.AddExpectation(expectation)
	return h
}

func (h *httpCheck) ExpectContent(expectedContent string) *httpCheck {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkContent(expectedContent, resp.Content)
	}
	h.AddExpectation(expectation)
	return h
}

func (h *httpCheck) ExpectBody(expectedBody string) *httpCheck {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkBody(expectedBody, resp.Content)
	}
	h.AddExpectation(expectation)
	return h
}

func (h *httpCheck) Auth(username string, password string) *httpCheck {
	h.username = username
	h.password = password
	return h
}

func (h *httpCheck) AuthFile(username string, passwordFile string) *httpCheck {
	h.username = username
	h.passwordFile = passwordFile
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
	expression := fmt.Sprintf(`(?is)<html[^>]*>.*?<head[^>]*>.*?<title[^>]*>[^<>]*%s[^<>]*</title>.*?</head>.*?</html>`, regexp.QuoteMeta(expectedTitle))
	logger.Tracef("title regexp: %s", expression)
	re := regexp.MustCompile(expression)
	if len(re.FindSubmatch(content)) > 0 {
		return nil
	}
	return fmt.Errorf("title %s not found", expectedTitle)
}

func checkStatusCode(expectedStatusCode int, statusCode int) error {
	if expectedStatusCode <= 0 {
		return nil
	}
	logger.Tracef("expectedStatusCode %d == statusCode %d", expectedStatusCode, statusCode)
	if expectedStatusCode != statusCode {
		return fmt.Errorf("wrong statuscode, expected %d got %d", expectedStatusCode, statusCode)
	}
	return nil
}

func get(executeRequest ExecuteRequest, url string, username string, password string) (*HttpResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(username) > 0 || len(password) > 0 {
		req.SetBasicAuth(username, password)
	}
	resp, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &HttpResponse{
		Content:    content,
		StatusCode: resp.StatusCode,
	}, nil
}
