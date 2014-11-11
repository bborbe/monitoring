package all

import (
	"runtime"
	"sync"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
)

var logger = log.DefaultLogger

type runnerAll struct {
}

func New() *runnerAll {
	return new(runnerAll)
}

func (r *runnerAll) Run(c configuration.Configuration) <-chan check.CheckResult {
	logger.Debug("run all checks")
	return Run(c.Checks())
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
