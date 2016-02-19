package all

import (
	"sync"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration "github.com/bborbe/monitoring/configuration"
	monitoring_node "github.com/bborbe/monitoring/node"
)

var logger = log.DefaultLogger

type runnerAll struct {
	maxConcurrency int
}

func New(maxConcurrency int) *runnerAll {
	r := new(runnerAll)
	r.maxConcurrency = maxConcurrency

	return r
}

func (r *runnerAll) Run(c monitoring_configuration.Configuration) <-chan monitoring_check.CheckResult {
	logger.Debug("run all checks")
	return Run(r.maxConcurrency, Checks(c))
}

func Run(maxConcurrency int, checks []monitoring_check.Check) <-chan monitoring_check.CheckResult {
	var wg sync.WaitGroup

	throttle := make(chan bool, maxConcurrency)

	resultChan := make(chan monitoring_check.CheckResult)

	for _, check := range checks {
		c := check
		wg.Add(1)
		go func() {
			throttle <- true
			resultChan <- c.Check()
			<-throttle
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
		logger.Debug("all checks finished")
	}()

	return resultChan
}

func Checks(configuration monitoring_configuration.Configuration) []monitoring_check.Check {
	list := make([]monitoring_check.Check, 0)
	list = addChecksToList(configuration.Nodes(), list)
	return list
}

func addChecksToList(nodes []monitoring_node.Node, checks []monitoring_check.Check) []monitoring_check.Check {
	if nodes != nil {
		for _, n := range nodes {
			if n.IsDisabled() {
				logger.Debugf("node %s disabled => skip", n.Check().Description())
				continue
			}
			if n.Check() != nil && !n.IsSilent() {
				checks = append(checks, n.Check())
			}
			checks = addChecksToList(n.Nodes(), checks)
		}
	}
	return checks
}
