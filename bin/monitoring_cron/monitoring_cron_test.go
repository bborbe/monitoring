package main

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/io"
)

func TestDoEmpty(t *testing.T) {
	writer := io.NewWriter()
	err := do(writer)
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
