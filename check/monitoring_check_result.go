package check

import "time"

type CheckResult interface {
	Success() bool
	Error() error
	Message() string
	Duration() time.Duration
}

type checkResult struct {
	success  bool
	error    error
	message  string
	duration time.Duration
}

func (c *checkResult) Success() bool {
	return c.success
}

func (c *checkResult) Error() error {
	return c.error
}

func (c *checkResult) Message() string {
	return c.message
}

func (c *checkResult) Duration() time.Duration {
	return c.duration
}

func NewCheckResult(chk Check, err error, duration time.Duration) CheckResult {
	if err != nil {
		return NewCheckResultFail(chk.Description(), err, duration)
	}
	return NewCheckResultSuccess(chk.Description(), duration)
}

func NewCheckResultSuccess(message string, duration time.Duration) CheckResult {
	r := new(checkResult)
	r.success = true
	r.message = message
	r.duration = duration
	r.error = nil
	return r
}

func NewCheckResultFail(message string, err error, duration time.Duration) CheckResult {
	r := new(checkResult)
	r.success = false
	r.message = message
	r.duration = duration
	r.error = err
	return r
}
