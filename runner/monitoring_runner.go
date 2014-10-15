package runner

import (
	"runtime"
	"sync"

	"github.com/bborbe/monitoring/check"
)

func Run(checks []check.Check) <-chan check.CheckResult {
	var wg sync.WaitGroup

	maxConcurrency := runtime.NumCPU() * 2
	throttle := make(chan bool, maxConcurrency)

	result := make(chan check.CheckResult)

	for _, check := range checks {
		c := check
		wg.Add(1)
		go func() {
			throttle <- true
			result <- c.Check()
			<-throttle
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}
