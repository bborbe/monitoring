package config

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/mailer"
)

func TestImplementsConfig(t *testing.T) {
	c := New()
	var i *mailer.Config
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
