package node

import (
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_check_dummy "github.com/bborbe/monitoring/check/dummy"
	monitoring_check_tcp "github.com/bborbe/monitoring/check/tcp"
)

func TestImplementsNode(t *testing.T) {
	c := New(monitoring_check_tcp.New("www.benjamin-borbe.de", 80))
	var i *Node
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithoutSubnode(t *testing.T) {
	var err error
	var node Node
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("succes"), "description")
	node = New(check)
	err = AssertThat(node, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(node.Nodes()), Is(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOneSubnode(t *testing.T) {
	var err error
	var node Node
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("succes"), "description")

	node = New(check, New(check))
	err = AssertThat(node, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(node.Nodes()), Is(1))
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsDisabledDefaultIsFalse(t *testing.T) {
	var err error
	var node Node
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("succes"), "description")

	node = New(check, New(check))
	err = AssertThat(node, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(node.IsDisabled(), Is(false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsSilentDefaultIsFalse(t *testing.T) {
	var err error
	var node Node
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("succes"), "description")

	node = New(check, New(check))
	err = AssertThat(node, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(node.IsSilent(), Is(false))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDisabledSetToTrue(t *testing.T) {
	var err error
	var node Node
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("succes"), "description")

	node = New(check, New(check)).Disabled(true)
	err = AssertThat(node, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(node.IsDisabled(), Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSilentSetToTrue(t *testing.T) {
	var err error
	var node Node
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("succes"), "description")

	node = New(check, New(check)).Silent(true)
	err = AssertThat(node, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(node.IsSilent(), Is(true))
	if err != nil {
		t.Fatal(err)
	}
}
