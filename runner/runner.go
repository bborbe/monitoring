package runner

import (
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_node "github.com/bborbe/monitoring/node"
)

type Runner interface {
	Run(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult
}
