package node

import (
	monitoring_check "github.com/bborbe/monitoring/check"
)

type Node interface {
	Check() monitoring_check.Check
	Nodes() []Node
	IsSilent() bool
	IsDisabled() bool
	Silent(silent bool) Node
	Disabled(disabled bool) Node
}

type node struct {
	check    monitoring_check.Check
	nodes    []Node
	silent   bool
	disabled bool
}

func New(check monitoring_check.Check, nodes ...Node) *node {
	n := new(node)
	n.check = check
	n.nodes = nodes
	return n
}

func (n *node) Check() monitoring_check.Check {
	return n.check
}

func (n *node) Nodes() []Node {
	return n.nodes
}

func (n *node) IsSilent() bool {
	return n.silent
}

func (n *node) IsDisabled() bool {
	return n.disabled
}

func (n *node) Silent(silent bool) Node {
	n.silent = silent
	return n
}

func (n *node) Disabled(disabled bool) Node {
	n.disabled = disabled
	return n
}
