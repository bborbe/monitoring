package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
)

type check struct {
	host string
}

var logger = log.DefaultLogger

func New(host string) *check {
	h := new(check)
	h.host = host
	return h
}

func (c *check) Check() monitoring_check.CheckResult {
	start := time.Now()
	ips, err := net.LookupHost(c.host)
	logger.Debugf("ips: %v", ips)
	return monitoring_check.NewCheckResult(c, err, time.Now().Sub(start))
}

func (c *check) Description() string {
	return fmt.Sprintf("dns check %s", c.host)
}
