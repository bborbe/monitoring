package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
)

type tcpCheck struct {
	host    string
	port    int
	timeout time.Duration
}

var logger = log.DefaultLogger

const (
	timeout = time.Duration(5 * time.Second)
	tries   = 3
)

func New(host string, port int) *tcpCheck {
	h := new(tcpCheck)
	h.host = host
	h.port = port
	return h
}

func (c *tcpCheck) Check() check.CheckResult {
	address := fmt.Sprintf("%s:%d", c.host, c.port)
	var err error
	for i := 0; i < tries; i++ {
		_, err = net.DialTimeout("tcp", address, timeout)
		logger.Debugf("tcp check on %s: %v", address, err)
		if err == nil {
			return check.NewCheckResult(c, err)
		}
	}
	return check.NewCheckResult(c, err)
}

func (c *tcpCheck) Description() string {
	return fmt.Sprintf("tcp check on %s:%d", c.host, c.port)
}
