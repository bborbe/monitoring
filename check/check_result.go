package check

type CheckResult interface {
	Success() bool
	Error() error
	Message() string
}

type checkResult struct {
	success bool
	error   error
	message string
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

func NewCheckResult(chk Check, err error) CheckResult {
	if err != nil {
		return NewCheckResultFail(chk.Description(), err)
	}
	return NewCheckResultSuccess(chk.Description())
}

func NewCheckResultSuccess(message string) CheckResult {
	r := new(checkResult)
	r.success = true
	r.message = message
	r.error = nil
	return r
}

func NewCheckResultFail(message string, err error) CheckResult {
	r := new(checkResult)
	r.success = false
	r.message = message
	r.error = err
	return r
}
