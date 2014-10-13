package configuration

import (
	"testing"
	. "github.com/bborbe/assert"
)

func TestImplementsConfiguration(t *testing.T) {
	c := New()
	var i *Configuration
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
