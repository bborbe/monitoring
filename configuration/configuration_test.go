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

func TestChecksFound(t *testing.T) {
	c := New()
	checks := c.Checks()
	err := AssertThat(checks, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(checks) > 0, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNodesFound(t *testing.T) {
	c := New()
	nodes := c.Nodes()
	err := AssertThat(nodes, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(nodes) > 0, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonParseFailed(t *testing.T) {
	err := AssertThat(checkBackupJson([]byte("")), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonSuccessParseAndNoMissingBackupFound(t *testing.T) {
	err := AssertThat(checkBackupJson([]byte("[]")), NilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonSuccessParseButMissingBackupFound(t *testing.T) {
	err := AssertThat(checkBackupJson([]byte("[{'a' : 'b'},{'b' : 'c'}]")), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
