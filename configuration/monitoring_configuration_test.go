package configuration

import (
	"testing"

	. "github.com/bborbe/assert"
	monitoring_check_http "github.com/bborbe/monitoring/check/http"
)

func TestImplementsConfiguration(t *testing.T) {
	c := New()
	var i *Configuration
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestNodesFound(t *testing.T) {
	c := New()
	nodes, err := c.Nodes()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(nodes, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(nodes) > 0, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonParseNilFailed(t *testing.T) {
	if err := AssertThat(checkBackupJson(&monitoring_check_http.HttpResponse{}), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonParseFailed(t *testing.T) {
	if err := AssertThat(checkBackupJson(&monitoring_check_http.HttpResponse{Content: []byte("")}), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonSuccessParseAndNoMissingBackupFound(t *testing.T) {
	if err := AssertThat(checkBackupJson(&monitoring_check_http.HttpResponse{Content: []byte("[]")}), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupJsonSuccessParseButMissingBackupFound(t *testing.T) {
	if err := AssertThat(checkBackupJson(&monitoring_check_http.HttpResponse{Content: []byte("[{'a' : 'b'},{'b' : 'c'}]")}), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
