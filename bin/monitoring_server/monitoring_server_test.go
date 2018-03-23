package main

import (
	"testing"

	"fmt"
	"time"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_node "github.com/bborbe/monitoring/node"
	"context"
)

func TestDoSendNoMail(t *testing.T) {
	counter := 0
	s := &server{
		runner: func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult {
			c := make(chan monitoring_check.CheckResult, 1)
			c <- monitoring_check.NewCheckResultSuccess("ok", time.Millisecond)
			close(c)
			return c
		},
		notify: func(sender string, recipient string, subject string, results []monitoring_check.CheckResult) error {
			counter++
			return nil
		},
		parseNodes: func(content string) ([]monitoring_node.Node, error) {
			return nil, nil
		},
		oneTime: true,
	}
	err := s.run(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestDoSendMail(t *testing.T) {
	counter := 0
	s := &server{
		runner: func(nodes []monitoring_node.Node) <-chan monitoring_check.CheckResult {
			c := make(chan monitoring_check.CheckResult, 1)
			c <- monitoring_check.NewCheckResultFail("ok", fmt.Errorf("foo"), time.Millisecond)
			close(c)
			return c
		},
		notify: func(sender string, recipient string, subject string, results []monitoring_check.CheckResult) error {
			counter++
			return nil
		},
		parseNodes: func(content string) ([]monitoring_node.Node, error) {
			return nil, nil
		},
		oneTime: true,
	}
	err := s.run(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
