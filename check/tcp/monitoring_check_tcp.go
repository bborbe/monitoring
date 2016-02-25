package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
)

type tcpCheck struct {
	host         string
	port         int
	timeout      time.Duration
	retryCounter int
}

var logger = log.DefaultLogger

const (
	DEFAULT_TIMEOUT       = time.Duration(5 * time.Second)
	DEFAULT_RETRY_COUNTER = 3
)

func New(host string, port int) *tcpCheck {
	h := new(tcpCheck)
	h.host = host
	h.port = port
	h.timeout = DEFAULT_TIMEOUT
	h.retryCounter = DEFAULT_RETRY_COUNTER
	return h
}

func (c *tcpCheck) Timeout(timeout time.Duration) *tcpCheck {
	c.timeout = timeout
	return c
}

func (c *tcpCheck) RetryCounter(retryCounter int) *tcpCheck {
	c.retryCounter = retryCounter
	return c
}

func (c *tcpCheck) Check() monitoring_check.CheckResult {
	start := time.Now()
	address := fmt.Sprintf("%s:%d", c.host, c.port)
	var err error
	for i := 0; i < c.retryCounter; i++ {
		_, err = net.DialTimeout("tcp", address, c.timeout)
		logger.Debugf("tcp check on %s: %v", address, err)
		if err == nil {
			return monitoring_check.NewCheckResult(c, err, time.Now().Sub(start))
		}
	}
	return monitoring_check.NewCheckResult(c, err, time.Now().Sub(start))
}

func (c *tcpCheck) Description() string {
	return fmt.Sprintf("tcp check on %s:%d", c.host, c.port)
}
