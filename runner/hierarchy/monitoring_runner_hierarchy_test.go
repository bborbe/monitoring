package hierarchy

import (
	"fmt"
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner "github.com/bborbe/monitoring/runner"
)

func TestImplementsRunner(t *testing.T) {
	c := New(1)
	var i *monitoring_runner.Runner
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	var err error
	c := NewCheck(monitoring_check.NewCheckResultSuccess("success"))
	nodes := make([]monitoring_node.Node, 0)
	nodes = append(nodes, monitoring_node.New(c))

	err = AssertThat(c.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

	resultChan := Run(1, nodes)
	result := <-resultChan

	err = AssertThat(c.counter, Is(1))
	if err != nil {
		t.Fatal(err)
	}

	err = AssertThat(result.Success(), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Error(), NilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Message(), Is("success"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunRecursive(t *testing.T) {
	var err error
	c := NewCheck(monitoring_check.NewCheckResultSuccess("success"))

	subnodes := make([]monitoring_node.Node, 0)
	subnodes = append(subnodes, monitoring_node.New(c))
	subnodes = append(subnodes, monitoring_node.New(c))

	nodes := make([]monitoring_node.Node, 0)
	nodes = append(nodes, monitoring_node.New(c, subnodes...))

	err = AssertThat(c.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

	resultChan := Run(1, nodes)
	<-resultChan
	<-resultChan
	<-resultChan

	err = AssertThat(c.counter, Is(3))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunRecursiveOnlyIfParentSuccess(t *testing.T) {
	var err error
	checkSuccess := NewCheck(monitoring_check.NewCheckResultSuccess("success"))
	checkFail := NewCheck(monitoring_check.NewCheckResultFail("fail", fmt.Errorf("foo")))

	subnodes := make([]monitoring_node.Node, 0)
	subnodes = append(subnodes, monitoring_node.New(checkSuccess))

	nodes := make([]monitoring_node.Node, 0)
	nodes = append(nodes, monitoring_node.New(checkFail, subnodes...))

	err = AssertThat(checkFail.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkSuccess.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

	resultChan := Run(1, nodes)
	<-resultChan

	err = AssertThat(checkFail.counter, Is(1))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkSuccess.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

}

type countCheck struct {
	counter int
	result  monitoring_check.CheckResult
}

func NewCheck(result monitoring_check.CheckResult) *countCheck {
	c := new(countCheck)
	c.counter = 0
	c.result = result
	return c
}

func (c *countCheck) Check() monitoring_check.CheckResult {
	c.counter++
	return c.result
}

func (c *countCheck) Description() string {
	return "description"
}
