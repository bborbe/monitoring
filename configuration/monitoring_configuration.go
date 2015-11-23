package configuration

import (
	"encoding/json"
	"fmt"

	"github.com/bborbe/monitoring/check/http"
	"github.com/bborbe/monitoring/check/tcp"
	"github.com/bborbe/monitoring/node"
)

type Configuration interface {
	Nodes() []node.Node
}

type configuration struct {
}

func New() Configuration {
	return new(configuration)
}

func (c *configuration) Nodes() []node.Node {
	list := make([]node.Node, 0)
	list = append(list, createNodeInternetAvaiable())
	return list
}

func createNodeInternetAvaiable() node.Node {
	return node.New(tcp.New("www.google.com", 80), createExternalNode(), createHmNode(), createRnNode(), createRaspVPN(), createRocketnewsVPN()).Silent(true)
}

func createExternalNode() node.Node {
	return node.New(http.New("http://benjaminborbe.zenfolio.com/").ExpectTitle("Zenfolio | Benjamin Borbe Fotografie"))
}

func createRnNode() node.Node {
	list := make([]node.Node, 0)

	list = append(list, node.New(tcp.New("144.76.187.199", 22)))
	list = append(list, node.New(tcp.New("144.76.187.200", 22)))
	list = append(list, node.New(tcp.New("144.76.187.199", 80)))
	list = append(list, node.New(tcp.New("144.76.187.200", 80)))
	list = append(list, node.New(tcp.New("144.76.187.199", 443)))
	list = append(list, node.New(tcp.New("144.76.187.200", 443)))

	list = append(list, node.New(http.New("http://www.benjamin-borbe.de/").ExpectTitle("Portfolio")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/").ExpectTitle("Portfolio")))
	list = append(list, node.New(http.New("http://www.benjaminborbe.de/").ExpectTitle("Portfolio")))
	list = append(list, node.New(http.New("https://www.benjaminborbe.de/").ExpectTitle("Portfolio")))

	list = append(list, node.New(http.New("http://www.benjamin-borbe.de/blog/").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/blog/").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("http://www.benjaminborbe.de/blog/").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("https://www.benjaminborbe.de/blog/").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("http://blog.benjamin-borbe.de/").ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/blog/").ExpectTitle("Benjamin Borbe Fotografie")))

	list = append(list, node.New(http.New("http://www.harteslicht.com/").ExpectTitle("Portfolio")))
	list = append(list, node.New(http.New("http://www.harteslicht.de/").ExpectTitle("Portfolio")))

	list = append(list, node.New(http.New("http://www.harteslicht.com/blog/").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, node.New(http.New("http://www.harteslicht.de/blog/").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, node.New(http.New("http://blog.harteslicht.com/").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, node.New(http.New("http://blog.harteslicht.de/").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))

	list = append(list, node.New(http.New("http://portfolio.benjamin-borbe.de/")))
	list = append(list, node.New(http.New("http://jana-und-ben.benjamin-borbe.de/").ExpectTitle("Portfolio")))
	list = append(list, node.New(http.New("http://jbf.benjamin-borbe.de/").ExpectTitle("Portfolio")))

	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/confluence/").ExpectTitle("Dashboard - Confluence")))
	list = append(list, node.New(http.New("http://confluence.benjamin-borbe.de/").ExpectTitle("Dashboard - Confluence")))

	list = append(list, node.New(http.New("http://portfolio.harteslicht.com/")))
	list = append(list, node.New(http.New("http://portfolio.harteslicht.de/")))

	list = append(list, node.New(http.New("http://kickstart.benjamin-borbe.de/").ExpectBody("ks.cfg")))
	list = append(list, node.New(http.New("http://ks.benjamin-borbe.de/").ExpectBody("ks.cfg")))

	list = append(list, node.New(http.New("http://slideshow.benjamin-borbe.de/").ExpectBody("go.html")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/slideshow/").ExpectBody("go.html")))

	list = append(list, node.New(http.New("http://jenkins.benjamin-borbe.de/").ExpectTitle("Dashboard [Jenkins]")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/jenkins/").ExpectTitle("Dashboard [Jenkins]")))

	list = append(list, node.New(http.New("http://ip.benjamin-borbe.de/")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/ip/")))

	list = append(list, node.New(http.New("http://backup.benjamin-borbe.de/").ExpectBody("Backup-Status")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/backup/").ExpectBody("Backup-Status")))

	list = append(list, node.New(http.New("http://booking.benjamin-borbe.de/status").ExpectContent("OK")))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/booking/status").ExpectContent("OK")))

	list = append(list, node.New(http.New("http://aptly.benjamin-borbe.de/").ExpectTitle(`Index of /`)))
	list = append(list, node.New(http.New("http://aptly.benjamin-borbe.de/api/version").Auth("api", "KYkobxZ6uvaGnYBG").ExpectContent(`{"Version":"0.9.5"}`)))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/aptly/").ExpectTitle(`Index of /`)))
	list = append(list, node.New(http.New("https://www.benjamin-borbe.de/aptly/api/version").Auth("api", "KYkobxZ6uvaGnYBG").ExpectContent(`{"Version":"0.9.5"}`)))

	list = append(list, createRnMailNode())

	return node.New(tcp.New("host.rocketsource.de", 22), list...)
}

func createRnMailNode() node.Node {
	list := make([]node.Node, 0)
	list = append(list, node.New(tcp.New("iredmail.mailfolder.org", 143)))
	list = append(list, node.New(tcp.New("iredmail.mailfolder.org", 993)))
	list = append(list, node.New(tcp.New("iredmail.mailfolder.org", 465)))
	return node.New(tcp.New("iredmail.mailfolder.org", 22), list...)
}

func createPnNode() node.Node {
	list := make([]node.Node, 0)
	var contentExpectation http.ContentExpectation
	contentExpectation = checkBackupJson
	list = append(list, node.New(http.New("http://backup.pn.benjamin-borbe.de:7777?status=false").AddExpectation(contentExpectation)))
	return node.New(tcp.New("backup.pn.benjamin-borbe.de", 7777), list...)
}

func createRaspVPN() node.Node {
	return node.New(tcp.New("10.30.0.1", 22), createPnNode()).Silent(true)
}

func createRocketnewsVPN() node.Node {
	return node.New(tcp.New("10.20.0.1", 22)).Silent(true)
}

func createHmNode() node.Node {
	list := make([]node.Node, 0)
	return node.New(tcp.New("home.benjamin-borbe.de", 443), list...)
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
