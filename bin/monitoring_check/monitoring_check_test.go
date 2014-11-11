package main

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/io"
	"github.com/bborbe/monitoring/runner/all"
)

func TestDoEmpty(t *testing.T) {
	writer := io.NewWriter()
	err := do(writer, all.New())
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
