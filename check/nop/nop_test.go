package nop

import (
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
)

func TestImplementsCheck(t *testing.T) {
	c := New("description")
	var i *monitoring_check.Check
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
