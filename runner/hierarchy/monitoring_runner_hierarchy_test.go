package hierarchy

import (
	"fmt"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/node"
	"github.com/bborbe/monitoring/runner"
)

func TestImplementsRunner(t *testing.T) {
	c := New()
	var i *runner.Runner
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	var err error
	c := NewCheck(check.NewCheckResultSuccess("success"))
	nodes := make([]node.Node, 0)
	nodes = append(nodes, node.New(c))

	err = AssertThat(c.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

	resultChan := Run(nodes)
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
	c := NewCheck(check.NewCheckResultSuccess("success"))

	subnodes := make([]node.Node, 0)
	subnodes = append(subnodes, node.New(c))
	subnodes = append(subnodes, node.New(c))

	nodes := make([]node.Node, 0)
	nodes = append(nodes, node.New(c, subnodes...))

	err = AssertThat(c.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

	resultChan := Run(nodes)
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
	checkSuccess := NewCheck(check.NewCheckResultSuccess("success"))
	checkFail := NewCheck(check.NewCheckResultFail("fail", fmt.Errorf("foo")))

	subnodes := make([]node.Node, 0)
	subnodes = append(subnodes, node.New(checkSuccess))

	nodes := make([]node.Node, 0)
	nodes = append(nodes, node.New(checkFail, subnodes...))

	err = AssertThat(checkFail.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkSuccess.counter, Is(0))
	if err != nil {
		t.Fatal(err)
	}

	resultChan := Run(nodes)
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
	result  check.CheckResult
}

func NewCheck(result check.CheckResult) *countCheck {
	c := new(countCheck)
	c.counter = 0
	c.result = result
	return c
}

func (c *countCheck) Check() check.CheckResult {
	c.counter++
	return c.result
}

func (c *countCheck) Description() string {
	return "description"
}
