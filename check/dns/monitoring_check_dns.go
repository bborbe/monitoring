package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
)

type dnsCheck struct {
	host string
}

var logger = log.DefaultLogger

func New(host string) *dnsCheck {
	h := new(dnsCheck)
	h.host = host
	return h
}

func (c *dnsCheck) Check() monitoring_check.CheckResult {
	start := time.Now()
	ips, err := net.LookupHost(c.host)
	logger.Debugf("ips: %v", ips)
	return monitoring_check.NewCheckResult(c, err, time.Now().Sub(start))
}

func (c *dnsCheck) Description() string {
	return fmt.Sprintf("dns check %s", c.host)
}
