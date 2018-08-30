package dns

import (
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
)

func TestImplementsCheck(t *testing.T) {
	c := New("www.example.com")
	var i *monitoring_check.Check
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDescription(t *testing.T) {
	c := New("www.example.com")
	err := AssertThat(c.Description(), Is("dns check www.example.com"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckSuccess(t *testing.T) {
	if testing.Short() {
		return
	}

	var err error
	c := New("www.benjamin-borbe.de")
	result := c.Check()
	err = AssertThat(result, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Success(), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Message(), Is("dns check www.benjamin-borbe.de"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Error(), NilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckFailure(t *testing.T) {
	if testing.Short() {
		return
	}
	c := New("notexistsing")
	result := c.Check()
	if err := AssertThat(result, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result.Success(), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result.Message(), Is("dns check notexistsing")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result.Error(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
