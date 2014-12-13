package hierarchy

import (
	"runtime"
	"sync"

	"github.com/bborbe/log"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/node"
)

type hierarchyRunner struct {
}

var logger = log.DefaultLogger

func New() *hierarchyRunner {
	return new(hierarchyRunner)
}

func (h *hierarchyRunner) Run(c configuration.Configuration) <-chan check.CheckResult {
	logger.Debug("run hierarchy checks")
	return Run(c.Nodes())
}

func Run(nodes []node.Node) <-chan check.CheckResult {
	resultChan := make(chan check.CheckResult)
	wg := new(sync.WaitGroup)

	maxConcurrency := runtime.NumCPU() * 2
	throttle := make(chan bool, maxConcurrency)

	wg.Add(1)
	go func() {
		exec(nodes, resultChan, wg, throttle)
		wg.Done()
	}()
	go func() {
		wg.Wait()
		close(resultChan)
		logger.Debug("all checks finished")
	}()

	return resultChan
}

func exec(nodes []node.Node, resultChan chan<- check.CheckResult, wg *sync.WaitGroup, throttle chan bool) {
	for _, n := range nodes {
		if n.IsDisabled() {
			logger.Debugf("node %s disabled => skip", n.Check().Description())
			continue
		}
		c := n.Check()
		ns := n.Nodes()
		isSilenet := n.IsSilent()
		wg.Add(1)
		go func() {
			throttle <- true
			result := c.Check()
			if !isSilenet {
				resultChan <- result
			}
			<-throttle
			if result.Success() && ns != nil {
				wg.Add(1)
				go func() {
					exec(ns, resultChan, wg, throttle)
					wg.Done()
				}()
			}
			wg.Done()
		}()
	}
}
