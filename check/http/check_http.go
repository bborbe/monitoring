package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bborbe/monitoring/check"
)

type httpCheck struct {
	url string
}

func New(url string) check.Check {
	h := new(httpCheck)
	h.url = url
	return h
}

func (h *httpCheck) Description() string {
	return fmt.Sprintf("http check on url %s", h.url)
}

func (h *httpCheck) Check() check.CheckResult {
	err := do(h.url)
	return check.NewCheckResult(h, err)
}

func do(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
