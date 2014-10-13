package http

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check"
)

func TestImplementsCheck(t *testing.T) {
	c := New()
	var i *check.Check
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
