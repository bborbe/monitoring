package configuration

import (
	"encoding/json"
	"fmt"

	monitoring_check_http "github.com/bborbe/monitoring/check/http"
	monitoring_check_tcp "github.com/bborbe/monitoring/check/tcp"
	monitoring_node "github.com/bborbe/monitoring/node"
)

type Configuration interface {
	Nodes() ([]monitoring_node.Node, error)
}

type configuration struct {
}

func New() Configuration {
	return new(configuration)
}

func (c *configuration) Nodes() ([]monitoring_node.Node, error) {
	list := make([]monitoring_node.Node, 0)
	list = append(list, createNodeInternetAvaiable())
	return list, nil
}

func createNodeInternetAvaiable() monitoring_node.Node {
	return monitoring_node.New(monitoring_check_tcp.New("www.google.com", 80),
		createExternalNode(),
		createHmNode(),
		createRnNode(),
		createRaspVPN(),
		createRocketnewsVPN(),
		createDevelNode()).Silent(true)
}

func createExternalNode() monitoring_node.Node {
	return monitoring_node.New(monitoring_check_http.New("http://benjaminborbe.zenfolio.com/").ExpectStatusCode(200).ExpectTitle("Zenfolio | Benjamin Borbe Fotografie"))
}

func createRnNode() monitoring_node.Node {
	list := make([]monitoring_node.Node, 0)

	list = append(list, monitoring_node.New(monitoring_check_tcp.New("144.76.187.199", 22)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("144.76.187.200", 22)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("144.76.187.199", 80)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("144.76.187.200", 80)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("144.76.187.199", 443)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("144.76.187.200", 443)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.benjaminborbe.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjaminborbe.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.benjamin-borbe.de/blog").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/blog").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.benjaminborbe.de/blog").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjaminborbe.de/blog").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.benjamin-borbe.de/blog/").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/blog/").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.benjaminborbe.de/blog/").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjaminborbe.de/blog/").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://blog.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjaminborbe.de/blog/").ExpectStatusCode(200).ExpectTitle("Benjamin Borbe Fotografie")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/wp-content").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/wp-content/").ExpectStatusCode(200)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/googlebd5f3e34a3e508a2.html").ExpectStatusCode(200).ExpectContent("google-site-verification: googlebd5f3e34a3e508a2.html")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.harteslicht.de/googlebd5f3e34a3e508a2.html").ExpectStatusCode(200).ExpectContent("google-site-verification: googlebd5f3e34a3e508a2.html")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.harteslicht.com/googlebd5f3e34a3e508a2.html").ExpectStatusCode(200).ExpectContent("google-site-verification: googlebd5f3e34a3e508a2.html")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.harteslicht.com/").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.harteslicht.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.harteslicht.com/blog/").ExpectStatusCode(200).ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.harteslicht.de/blog/").ExpectStatusCode(200).ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://blog.harteslicht.com/").ExpectStatusCode(200).ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://blog.harteslicht.de/").ExpectStatusCode(200).ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht.")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://portfolio.benjamin-borbe.de/")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://jana-und-ben.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/jana-und-ben").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/jana-und-ben/").ExpectStatusCode(200).ExpectTitle("Portfolio")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://jbf.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Portfolio")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://confluence.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Dashboard - Confluence")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/confluence").ExpectStatusCode(200).ExpectTitle("Dashboard - Confluence")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/confluence/").ExpectStatusCode(200).ExpectTitle("Dashboard - Confluence")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://portfolio.harteslicht.com/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://portfolio.harteslicht.de/").ExpectStatusCode(200)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://kickstart.benjamin-borbe.de/").ExpectStatusCode(200).ExpectBody("ks.cfg")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/kickstart").ExpectStatusCode(200).ExpectBody("ks.cfg")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/kickstart/").ExpectStatusCode(200).ExpectBody("ks.cfg")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://ks.benjamin-borbe.de/").ExpectStatusCode(200).ExpectBody("ks.cfg")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://slideshow.benjamin-borbe.de/").ExpectStatusCode(200).ExpectBody("go.html")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/slideshow").ExpectStatusCode(200).ExpectBody("go.html")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/slideshow/").ExpectStatusCode(200).ExpectBody("go.html")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://jenkins.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle("Dashboard [Jenkins]")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/jenkins").ExpectStatusCode(200).ExpectTitle("Dashboard [Jenkins]")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/jenkins/").ExpectStatusCode(200).ExpectTitle("Dashboard [Jenkins]")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://ip.benjamin-borbe.de/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/ip").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/ip/").ExpectStatusCode(200)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://password.benjamin-borbe.de/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/password").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/password/").ExpectStatusCode(200)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://rocketnews.de/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.rocketnews.de/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://rocketsource.de/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://www.rocketsource.de/").ExpectStatusCode(200)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://backup.benjamin-borbe.de/").ExpectStatusCode(200).ExpectBody("Backup-Status")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/backup").ExpectStatusCode(200).ExpectBody("Backup-Status")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/backup/").ExpectStatusCode(200).ExpectBody("Backup-Status")))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://booking.benjamin-borbe.de/status").ExpectStatusCode(200).ExpectStatusCode(200).ExpectContent("OK")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/booking/status").ExpectStatusCode(200).ExpectStatusCode(200).ExpectContent("OK")))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://booking.benjamin-borbe.de/").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/booking").ExpectStatusCode(200)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/booking/").ExpectStatusCode(200)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("http://aptly.benjamin-borbe.de/").ExpectStatusCode(200).ExpectTitle(`Index of /`)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/aptly").ExpectStatusCode(200).ExpectTitle(`Index of /`)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/aptly/").ExpectStatusCode(200).ExpectTitle(`Index of /`)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://aptly.benjamin-borbe.de/api/version").ExpectStatusCode(200).AuthFile("api", "/etc/aptly_api_password").ExpectContent(`{"Version":"0.9.6"}`)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/aptly/api/version").ExpectStatusCode(200).AuthFile("api", "/etc/aptly_api_password").ExpectContent(`{"Version":"0.9.6"}`)))

	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/webdav").ExpectStatusCode(401)))
	list = append(list, monitoring_node.New(monitoring_check_http.New("https://www.benjamin-borbe.de/webdav/").ExpectStatusCode(401)))

	list = append(list, createRnMailNode())

	return monitoring_node.New(monitoring_check_tcp.New("host.rocketsource.de", 22), list...)
}

func createRnMailNode() monitoring_node.Node {
	list := make([]monitoring_node.Node, 0)
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("iredmail.mailfolder.org", 143)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("iredmail.mailfolder.org", 993)))
	list = append(list, monitoring_node.New(monitoring_check_tcp.New("iredmail.mailfolder.org", 465)))
	return monitoring_node.New(monitoring_check_tcp.New("iredmail.mailfolder.org", 22), list...)
}

func createPnNode() monitoring_node.Node {
	list := make([]monitoring_node.Node, 0)
	var contentExpectation monitoring_check_http.Expectation
	contentExpectation = checkBackupJson
	list = append(list, monitoring_node.New(monitoring_check_http.New("http://backup.pn.benjamin-borbe.de:7777?status=false").ExpectStatusCode(200).AddExpectation(contentExpectation)))
	return monitoring_node.New(monitoring_check_tcp.New("backup.pn.benjamin-borbe.de", 7777), list...)
}

func createDevelNode() monitoring_node.Node {
	return monitoring_node.New(monitoring_check_tcp.New("bborbe.devel.lf.seibert-media.net", 22),
		monitoring_node.New(monitoring_check_tcp.New("bborbe.devel.lf.seibert-media.net", 80),
			monitoring_node.New(monitoring_check_http.New("http://bborbe.devel.lf.seibert-media.net/").ExpectStatusCode(200)),
			monitoring_node.New(monitoring_check_http.New("http://bborbe.devel.lf.seibert-media.net/dash/").ExpectStatusCode(200).ExpectTitle("Dash")),
			monitoring_node.New(monitoring_check_http.New("http://bborbe.devel.lf.seibert-media.net/collaborate/").ExpectStatusCode(200).ExpectContent("collaboratetheme")),
		),
		monitoring_node.New(monitoring_check_tcp.New("bborbe.devel.lf.seibert-media.net", 443),
			monitoring_node.New(monitoring_check_http.New("https://bborbe.devel.lf.seibert-media.net/").ExpectStatusCode(200)),
			monitoring_node.New(monitoring_check_http.New("https://bborbe.devel.lf.seibert-media.net/dash/").ExpectStatusCode(200).ExpectTitle("Dash")),
		),
	).Silent(true)
}

func createRaspVPN() monitoring_node.Node {
	return monitoring_node.New(monitoring_check_tcp.New("10.30.0.1", 22), createPnNode()).Silent(true)
}

func createRocketnewsVPN() monitoring_node.Node {
	return monitoring_node.New(monitoring_check_tcp.New("10.20.0.1", 22)).Silent(true)
}

func createHmNode() monitoring_node.Node {
	list := make([]monitoring_node.Node, 0)
	return monitoring_node.New(monitoring_check_tcp.New("home.benjamin-borbe.de", 443), list...)
}

func checkBackupJson(resp *monitoring_check_http.HttpResponse) error {
	var data []interface{}
	err := json.Unmarshal(resp.Content, &data)
	if err != nil {
		return fmt.Errorf("parse json failed")
	}
	if len(data) > 0 {
		return fmt.Errorf("found false backups")
	}
	return nil
}
