package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
)

func TestDoEmpty(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, monitoring_runner_all.New(1), NewConfiguration())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.String(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(writer.String()) > 0, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

type configurationDummy struct{}

func NewConfiguration() *configurationDummy {
	return new(configurationDummy)
}

func (c *configurationDummy) Checks() []monitoring_check.Check {
	return make([]monitoring_check.Check, 0)
}

func (c *configurationDummy) Nodes() []monitoring_node.Node {
	return make([]monitoring_node.Node, 0)
}
