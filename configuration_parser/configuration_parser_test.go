package configuration_parser

import (
	"testing"

	. "github.com/bborbe/assert"
	"reflect"
)

func TestImplementsConfigurationParser(t *testing.T) {
	c := New()
	var i *ConfigurationParser
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestParseEmptyConfigurationReturnError(t *testing.T) {
	c := New()
	_, err := c.ParseConfiguration([]byte(``))
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestParseEmptyNodes(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNode(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNodeSilentTrue(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node silent="true"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes[0].IsSilent(), Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNodeSilentFalse(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node silent="false"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes[0].IsSilent(), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNodeSilentNotSet(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes[0].IsSilent(), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNodeDisabledTrue(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node disabled="true"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes[0].IsDisabled(), Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNodeDisabledFalse(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node disabled="false"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes[0].IsDisabled(), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestParseOneNodeDisabledNotSet(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes[0].IsDisabled(), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestParseTcpCheck(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node check="tcp"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(reflect.TypeOf(nodes[0].Check()).String(), Is("*tcp.tcpCheck")); err != nil {
		t.Fatal(err)
	}
}

func TestParseHttpCheck(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node check="http"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(reflect.TypeOf(nodes[0].Check()).String(), Is("*http.httpCheck")); err != nil {
		t.Fatal(err)
	}
}

func TestParseDummyCheck(t *testing.T) {
	c := New()
	nodes, err := c.ParseConfiguration([]byte(`<nodes><node check="dummy"></node></nodes>`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(reflect.TypeOf(nodes[0].Check()).String(), Is("*dummy.checkDummy")); err != nil {
		t.Fatal(err)
	}
}