package configuration

import (
	"encoding/json"
	"fmt"

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
	list = append(list, createNodeInternetAvaiable())
	return list
}

func createNodeInternetAvaiable() node.Node {
	return node.New(tcp.New("www.google.com", 80), createNodeRocketsourceAvaiable())
}

func createNodeRocketsourceAvaiable() node.Node {
	list := make([]node.Node, 0)

	list = append(list, node.New(tcp.New("144.76.187.199", 22)))
	list = append(list, node.New(tcp.New("144.76.187.200", 22)))
	list = append(list, node.New(tcp.New("144.76.187.199", 80)))
	list = append(list, node.New(tcp.New("144.76.187.200", 80)))
	list = append(list, node.New(tcp.New("144.76.187.199", 443)))
	list = append(list, node.New(tcp.New("144.76.187.200", 443)))

	list = append(list, node.New(http.New("http://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("http://www.benjaminborbe.de").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/confluence/").ExpectTitle("Dashboard - Confluence")))

	list = append(list, node.New(http.New("http://www.harteslicht.de").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, node.New(http.New("http://www.harteslicht.com").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))

	list = append(list, node.New(http.New("http://jenkins.benjamin-borbe.de").ExpectTitle("Dashboard [Jenkins]")))
	list = append(list, node.New(http.New("http://kickstart.benjamin-borbe.de").ExpectBody("ks.cfg")))

	list = append(list, node.New(http.New("http://ip.benjamin-borbe.de")))
	list = append(list, node.New(http.New("http://slideshow.benjamin-borbe.de").ExpectBody("go.html")))
	list = append(list, node.New(http.New("http://apt.benjamin-borbe.de/bborbe-unstable/Sources").ExpectContent("bborbe-unstable")))
	list = append(list, node.New(http.New("http://blog.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie")))

	list = append(list, createBackupStatusNode())

	return node.New(tcp.New("host.rocketsource.de", 22), list...)
}

func createBackupStatusNode() node.Node {
	list := make([]node.Node, 0)
	list = append(list, node.New(http.New("http://backup.pn.benjamin-borbe.de:7777?status=false").AddExpectation(checkBackupJson)))
	return node.New(tcp.New("backup.pn.benjamin-borbe.de", 7777), list...)
}

func checkBackupJson(content []byte) error {
	var data []interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		return fmt.Errorf("parse json failed")
	}
	if len(data) > 0 {
		return fmt.Errorf("found false backups")
	}
	return nil
}
