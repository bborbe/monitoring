package node

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/monitoring/check/tcp"
)

func TestImplementsNode(t *testing.T) {
	c := New(tcp.New("www.benjamin-borbe.de", 80), make([]Node, 0))
	var i *Node
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
