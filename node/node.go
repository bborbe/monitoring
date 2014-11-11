package node

import (
	"github.com/bborbe/monitoring/check"
)

type Node interface {
	Check() check.Check
	Nodes() []Node
}

type node struct {
	check check.Check
	nodes []Node
}

func New(c check.Check, nodes []Node) *node {
	n := new(node)
	n.check = c
	n.nodes = nodes
	return n
}

func (n *node) Check() check.Check {
	return n.check
}

func (n *node) Nodes() []Node {
	return n.nodes
}
