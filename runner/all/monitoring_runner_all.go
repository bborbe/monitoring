package all

import (
	"runtime"
	"sync"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/node"
)

var logger = log.DefaultLogger

type runnerAll struct {
}

func New() *runnerAll {
	return new(runnerAll)
}

func (r *runnerAll) Run(c configuration.Configuration) <-chan check.CheckResult {
	logger.Debug("run all checks")
	return Run(Checks(c))
}

func Run(checks []check.Check) <-chan check.CheckResult {
	var wg sync.WaitGroup

	maxConcurrency := runtime.NumCPU() * 2
	throttle := make(chan bool, maxConcurrency)

	resultChan := make(chan check.CheckResult)

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

func Checks(c configuration.Configuration) []check.Check {
	list := make([]check.Check, 0)
	list = addChecksToList(c.Nodes(), list)
	return list
}

func addChecksToList(nodes []node.Node, checks []check.Check) []check.Check {
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
