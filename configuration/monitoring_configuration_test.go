package configuration

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check/http"
)

func TestImplementsConfiguration(t *testing.T) {
	c := New()
	var i *Configuration
	err := AssertThat(c, Implements(i))
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

func TestCheckBackupJsonParseNilFailed(t *testing.T) {
	if err := AssertThat(checkBackupJson(&http.HttpResponse{}), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonParseFailed(t *testing.T) {
	if err := AssertThat(checkBackupJson(&http.HttpResponse{Content: []byte("")}), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonSuccessParseAndNoMissingBackupFound(t *testing.T) {
	if err := AssertThat(checkBackupJson(&http.HttpResponse{Content: []byte("[]")}), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonSuccessParseButMissingBackupFound(t *testing.T) {
	if err := AssertThat(checkBackupJson(&http.HttpResponse{Content: []byte("[{'a' : 'b'},{'b' : 'c'}]")}), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
