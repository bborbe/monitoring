package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/node"
	"github.com/bborbe/monitoring/runner/all"
)

func TestDoEmpty(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, all.New(1), NewConfiguration())
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

func (c *configurationDummy) Checks() []check.Check {
	return make([]check.Check, 0)
}

func (c *configurationDummy) Nodes() []node.Node {
	return make([]node.Node, 0)
}
