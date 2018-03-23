package mailer

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsMail(t *testing.T) {
	c := New(nil)
	var i *Mailer
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuoteString(t *testing.T) {
	if err := AssertThat(QuoteString("hello world"), Is("hello world")); err != nil {
		t.Fatal(err)
	}
}

func TestQuoteStringSpecial(t *testing.T) {
	if err := AssertThat(QuoteString("hello wörld"), Is("hello w=C3=B6rld")); err != nil {
		t.Fatal(err)
	}
}
