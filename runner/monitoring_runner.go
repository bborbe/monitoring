package runner

import (
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
)

type Runner interface {
	Run(c configuration.Configuration) <-chan check.CheckResult
}
