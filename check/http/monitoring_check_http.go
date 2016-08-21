package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"time"

	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/redirect_follower"
	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
)

const (
	DEFAULT_TIMEOUT       = 30 * time.Second
	DEFAULT_RETRY_COUNTER = 3
	USERAGENT             = "Monitoring"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type Expectation func(httpResponse *HttpResponse) error

type check struct {
	url          string
	username     string
	password     string
	passwordFile string
	expectations []Expectation
	timeout      time.Duration
	retryCounter int
}

type HttpResponse struct {
	Content    []byte
	StatusCode int
}

var logger = log.DefaultLogger

func New(url string) *check {
	h := new(check)
	h.url = url
	h.timeout = DEFAULT_TIMEOUT
	h.retryCounter = DEFAULT_RETRY_COUNTER
	return h
}

func (c *check) Description() string {
	return fmt.Sprintf("http check on url %s", c.url)
}

func (c *check) executeRequest() ExecuteRequest {
	builder := http_client_builder.New().WithoutProxy().WithTimeout(c.timeout)
	redirectFollower := redirect_follower.New(builder.BuildRoundTripper().RoundTrip)
	return redirectFollower.ExecuteRequestAndFollow
}

func (c *check) RetryCounter(retryCounter int) *check {
	c.retryCounter = retryCounter
	return c
}

func (h *check) Timeout(timeout time.Duration) *check {
	h.timeout = timeout
	return h
}

func (c *check) Check() monitoring_check.CheckResult {
	start := time.Now()
	var err error
	for i := 0; i < c.retryCounter; i++ {
		err = c.check()
		if err == nil {
			break
		}
	}
	return monitoring_check.NewCheckResult(c, err, time.Now().Sub(start))
}

func (c *check) check() error {
	if len(c.password) == 0 && len(c.passwordFile) > 0 {
		logger.Debugf("read password from file %s", c.passwordFile)
		password, err := ioutil.ReadFile(c.passwordFile)
		if err != nil {
			logger.Debugf("read password file failed %s: %v", c.passwordFile, err)
			return err
		}
		c.password = strings.TrimSpace(string(password))
	}
	httpResponse, err := get(c.executeRequest(), c.url, c.username, c.password)
	if err != nil {
		logger.Debugf("fetch url failed %s: %v", c.url, err)
		return err
	}
	for _, expectation := range c.expectations {
		if err = expectation(httpResponse); err != nil {
			return err
		}
	}
	return nil
}

func (c *check) AddExpectation(expectation Expectation) *check {
	c.expectations = append(c.expectations, expectation)
	return c
}

func (c *check) ExpectTitle(expectedTitle string) *check {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkTitle(expectedTitle, resp.Content)
	}
	c.AddExpectation(expectation)
	return c
}

func (c *check) ExpectStatusCode(expectedStatusCode int) *check {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkStatusCode(expectedStatusCode, resp.StatusCode)
	}
	c.AddExpectation(expectation)
	return c
}

func (c *check) ExpectContent(expectedContent string) *check {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkContent(expectedContent, resp.Content)
	}
	c.AddExpectation(expectation)
	return c
}

func (c *check) ExpectBody(expectedBody string) *check {
	var expectation Expectation
	expectation = func(resp *HttpResponse) error {
		return checkBody(expectedBody, resp.Content)
	}
	c.AddExpectation(expectation)
	return c
}

func (c *check) Auth(username string, password string) *check {
	c.username = username
	c.password = password
	return c
}

func (c *check) AuthFile(username string, passwordFile string) *check {
	c.username = username
	c.passwordFile = passwordFile
	return c
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
	req.Header.Set("User-Agent", USERAGENT)
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
