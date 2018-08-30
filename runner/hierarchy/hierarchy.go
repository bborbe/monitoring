package hierarchy

import (
	"sync"

	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_node "github.com/bborbe/monitoring/node"
	"github.com/golang/glog"
)

type hierarchyRunner struct {
	maxConcurrency int
}

func New(maxConcurrency int) *hierarchyRunner {
	h := new(hierarchyRunner)
	h.maxConcurrency = maxConcurrency
	return h
}

func (h *hierarchyRunner) Run(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult {
	glog.V(2).Info("run hierarchy checks")
	return Run(h.maxConcurrency, nodes)
}

func Run(maxConcurrency int, nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult {
	resultChan := make(chan monitoring_check.CheckResult)
	wg := new(sync.WaitGroup)

	throttle := make(chan bool, maxConcurrency)

	wg.Add(1)
	go func() {
		exec(nodes, resultChan, wg, throttle)
		wg.Done()
	}()
	go func() {
		wg.Wait()
		close(resultChan)
		glog.V(2).Info("all checks finished")
	}()

	return resultChan
}

func exec(nodes []monitoring_node.Node, resultChan chan<- monitoring_check.CheckResult, wg *sync.WaitGroup, throttle chan bool) {
	for _, n := range nodes {
		if n.IsDisabled() {
			glog.V(2).Infof("node %s disabled => skip", n.Check().Description())
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
