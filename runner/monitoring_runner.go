package runner

import (
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration "github.com/bborbe/monitoring/configuration"
)

type Runner interface {
	Run(configuration monitoring_configuration.Configuration) <-chan monitoring_check.CheckResult
}
