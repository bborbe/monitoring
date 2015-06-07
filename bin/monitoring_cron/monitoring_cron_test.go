package main

import (
	"testing"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/configuration"
	"github.com/bborbe/monitoring/node"
	"github.com/bborbe/monitoring/runner/all"
)

func TestDoEmpty(t *testing.T) {
	writer := io_mock.NewWriter()
	r := all.New()
	err := do(writer, r, NewConfigurationDummy(make([]check.Check, 0), make([]node.Node, 0)), nil)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.Content(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(writer.Content()) > 0, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

type configurationDummy struct {
	checks []check.Check
	nodes  []node.Node
}

func NewConfigurationDummy(checks []check.Check, nodes []node.Node) configuration.Configuration {
	c := new(configurationDummy)
	c.checks = checks
	c.nodes = nodes
	return c
}

func (c *configurationDummy) Checks() []check.Check { return c.checks }
func (c *configurationDummy) Nodes() []node.Node    { return c.nodes }
