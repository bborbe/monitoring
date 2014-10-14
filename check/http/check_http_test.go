package http

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check"
)

func TestImplementsCheck(t *testing.T) {
	c := New("http://www.example.com")
	var i *check.Check
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckTitle(t *testing.T) {
	var err error
	err = AssertThat(checkTitle("test", []byte("<html><head><title>test</title></head></html>")), NilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkTitle("test", []byte("<html><head><title>foobar</title></head></html>")), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkTitle("test", []byte("<html>asdf<head>asdf<title>test</title>asdf</head>asdf</html>")), NilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkTitle("test", []byte("<html id=\"a\"><head id=\"b\"><title id=\"c\">test</title></head></html>")), NilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkTitle("test", []byte("<html>\n<head><title>test</title></head></html>")), NilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(checkTitle("test", []byte("<html><head><TiTle>test</title></head></html>")), NilValue())
	if err != nil {
		t.Fatal(err)
	}
}
