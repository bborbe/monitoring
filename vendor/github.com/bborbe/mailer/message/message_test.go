package message

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/mailer"
)

func TestImplementsMessage(t *testing.T) {
	c := New()
	var i *mailer.Message
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
