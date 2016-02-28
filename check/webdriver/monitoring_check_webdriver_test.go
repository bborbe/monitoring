package webdriver

import (
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
)

func TestImplementsCheck(t *testing.T) {
	c := New(nil, "http://www.example.com")
	var i *monitoring_check.Check
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDescription(t *testing.T) {
	c := New(nil, "http://www.example.com")
	err := AssertThat(c.Description(), Is("webdriver check on url http://www.example.com"))
	if err != nil {
		t.Fatal(err)
	}
}
