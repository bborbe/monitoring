package tcp

import (
	"testing"
	"time"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
)

func TestImplementsCheck(t *testing.T) {
	c := New("www.example.com", 80)
	var i *monitoring_check.Check
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDescription(t *testing.T) {
	c := New("www.example.com", 80)
	err := AssertThat(c.Description(), Is("tcp check on www.example.com:80"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckSuccess(t *testing.T) {
	if testing.Short() {
		return
	}

	var err error
	c := New("www.benjamin-borbe.de", 80)
	result := c.Check()
	err = AssertThat(result, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Success(), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Message(), Is("tcp check on www.benjamin-borbe.de:80"))
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

	var err error
	c := New("www.benjamin-borbe.de", 82)
	result := c.Check()
	err = AssertThat(result, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Success(), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Message(), Is("tcp check on www.benjamin-borbe.de:81"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result.Error(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultTimeout(t *testing.T) {
	c := New("www.benjamin-borbe.de", 81)
	if err := AssertThat(c.timeout, Is(5*time.Second)); err != nil {
		t.Fatal(err)
	}
}

func TestSetTimeout(t *testing.T) {
	c := New("www.benjamin-borbe.de", 81)
	c.Timeout(30 * time.Second)
	if err := AssertThat(c.timeout, Is(30*time.Second)); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultRetryCounter(t *testing.T) {
	c := New("www.benjamin-borbe.de", 81)
	if err := AssertThat(c.retryCounter, Is(3)); err != nil {
		t.Fatal(err)
	}
}

func TestRetryCounter(t *testing.T) {
	c := New("www.benjamin-borbe.de", 81)
	c.RetryCounter(25)
	if err := AssertThat(c.retryCounter, Is(25)); err != nil {
		t.Fatal(err)
	}
}
