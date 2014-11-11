package configuration

import (
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/check/http"
	"github.com/bborbe/monitoring/check/tcp"
	"github.com/bborbe/monitoring/node"
)

type Node interface {
	Check() check.Check
	Nodes() []Node
}

type Configuration interface {
	Checks() []check.Check
	Nodes() []node.Node
}

type configuration struct {
}

func New() Configuration {
	return new(configuration)
}

func (c *configuration) Checks() []check.Check {
	list := make([]check.Check, 0)
	list = addChecksToList(c.Nodes(), list)
	return list
}

func addChecksToList(nodes []node.Node, checks []check.Check) []check.Check {
	if nodes != nil {
		for _, n := range nodes {
			if n.Check() != nil {
				checks = append(checks, n.Check())
			}
			checks = addChecksToList(n.Nodes(), checks)
		}
	}
	return checks
}

func (c *configuration) Nodes() []node.Node {
	list := make([]node.Node, 0)
	list = append(list, node.New(tcp.New("144.76.187.199", 22), nil))
	list = append(list, node.New(tcp.New("144.76.187.200", 22), nil))
	list = append(list, node.New(tcp.New("144.76.187.199", 80), nil))
	list = append(list, node.New(tcp.New("144.76.187.200", 80), nil))
	list = append(list, node.New(tcp.New("144.76.187.199", 443), nil))
	list = append(list, node.New(tcp.New("144.76.187.200", 443), nil))

	list = append(list, node.New(http.New("http://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"), nil))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"), nil))
	list = append(list, node.New(http.New("http://www.benjaminborbe.de").ExpectTitle("Benjamin Borbe Fotografie"), nil))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/confluence/").ExpectTitle("Dashboard - Confluence"), nil))

	list = append(list, node.New(http.New("http://www.harteslicht.de").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht."), nil))
	list = append(list, node.New(http.New("http://www.harteslicht.com").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht."), nil))

	list = append(list, node.New(http.New("http://jenkins.benjamin-borbe.de").ExpectTitle("Dashboard [Jenkins]"), nil))
	list = append(list, node.New(http.New("http://kickstart.benjamin-borbe.de").ExpectBody("ks.cfg"), nil))

	list = append(list, node.New(http.New("http://ip.benjamin-borbe.de"), nil))
	list = append(list, node.New(http.New("http://slideshow.benjamin-borbe.de").ExpectBody("go.html"), nil))
	list = append(list, node.New(http.New("http://apt.benjamin-borbe.de/bborbe-unstable/Sources").ExpectContent("bborbe-unstable"), nil))
	list = append(list, node.New(http.New("http://blog.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"), nil))

	return list
}
