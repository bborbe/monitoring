package mock

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/mailer"
)

func TestImplementsMail(t *testing.T) {
	c := New()
	var i *mailer.Mailer
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
