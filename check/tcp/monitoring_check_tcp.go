package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
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

func (c *tcpCheck) Check() monitoring_check.CheckResult {
	start := time.Now()
	address := fmt.Sprintf("%s:%d", c.host, c.port)
	var err error
	for i := 0; i < tries; i++ {
		_, err = net.DialTimeout("tcp", address, timeout)
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
