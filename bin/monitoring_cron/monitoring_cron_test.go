package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_configuration "github.com/bborbe/monitoring/configuration"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
)

func TestDoEmpty(t *testing.T) {
	writer := bytes.NewBufferString("")
	r := monitoring_runner_all.New(1)
	err := do(writer, r, NewConfigurationDummy(make([]monitoring_check.Check, 0), make([]monitoring_node.Node, 0)), nil)
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

type configurationDummy struct {
	checks []monitoring_check.Check
	nodes  []monitoring_node.Node
}

func NewConfigurationDummy(checks []monitoring_check.Check, nodes []monitoring_node.Node) monitoring_configuration.Configuration {
	c := new(configurationDummy)
	c.checks = checks
	c.nodes = nodes
	return c
}

func (c *configurationDummy) Checks() []monitoring_check.Check { return c.checks }
func (c *configurationDummy) Nodes() []monitoring_node.Node    { return c.nodes }
