package tcp

import (
	"fmt"
	"net"
	"time"

	monitoring_check "github.com/bborbe/monitoring/check"
	"github.com/golang/glog"
)

type check struct {
	host         string
	port         int
	timeout      time.Duration
	retryCounter int
}

const (
	DEFAULT_TIMEOUT       = time.Duration(5 * time.Second)
	DEFAULT_RETRY_COUNTER = 3
)

func New(host string, port int) *check {
	h := new(check)
	h.host = host
	h.port = port
	h.timeout = DEFAULT_TIMEOUT
	h.retryCounter = DEFAULT_RETRY_COUNTER
	return h
}

func (c *check) Timeout(timeout time.Duration) *check {
	c.timeout = timeout
	return c
}

func (c *check) RetryCounter(retryCounter int) *check {
	c.retryCounter = retryCounter
	return c
}

func (h *check) Check() monitoring_check.CheckResult {
	start := time.Now()
	var err error
	for i := 0; i < h.retryCounter; i++ {
		err = h.check()
		if err == nil {
			break
		}
	}
	return monitoring_check.NewCheckResult(h, err, time.Now().Sub(start))
}

func (c *check) check() error {
	address := fmt.Sprintf("%s:%d", c.host, c.port)
	var err error
	_, err = net.DialTimeout("tcp", address, c.timeout)
	glog.V(2).Infof("tcp check on %s: %v", address, err)
	return err
}

func (c *check) Description() string {
	return fmt.Sprintf("tcp check on %s:%d", c.host, c.port)
}
