package main

import (
	"bytes"
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_node "github.com/bborbe/monitoring/node"
)

func TestDoEmpty(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult {
		return nil
	}, func(sender string, recipient string, subject string, results []monitoring_check.CheckResult) error {
		return nil
	}, func(content []byte) ([]monitoring_node.Node, error) {
		return nil, nil
	}, "", "/tmp/lock", "a", "b", "c")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
