package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
	monitoring_node "github.com/bborbe/monitoring/node"
	monitoring_runner_all "github.com/bborbe/monitoring/runner/all"
)

func TestDoEmpty(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, monitoring_runner_all.New(1), func() ([]monitoring_node.Node, error) {
		return make([]monitoring_node.Node, 0), nil
	})
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(writer.String(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(writer.String()) > 0, Is(true)); err != nil {
		t.Fatal(err)
	}
}
