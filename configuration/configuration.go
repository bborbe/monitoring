package configuration

import (
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/check/http"
	"github.com/bborbe/monitoring/check/tcp"
)

type Configuration interface {
	Checks() []check.Check
}

type configuration struct {
}

func New() Configuration {
	return new(configuration)
}

func (c *configuration) Checks() []check.Check {
	list := make([]check.Check, 0)
	list = append(list, tcp.New("144.76.187.199", 22))
	list = append(list, tcp.New("144.76.187.200", 22))
	list = append(list, tcp.New("144.76.187.199", 80))
	list = append(list, tcp.New("144.76.187.200", 80))
	list = append(list, tcp.New("144.76.187.199", 443))
	list = append(list, tcp.New("144.76.187.200", 443))
	list = append(list, http.New("http://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	list = append(list, http.New("https://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	list = append(list, http.New("http://www.benjaminborbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	list = append(list, http.New("http://jenkins.benjamin-borbe.de").ExpectTitle("Dashboard [Jenkins]"))
	list = append(list, http.New("http://www.harteslicht.de").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht."))
	list = append(list, http.New("http://www.harteslicht.com").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht."))
	list = append(list, http.New("http://kickstart.benjamin-borbe.de").ExpectBody("ks.cfg"))
	list = append(list, http.New("http://ip.benjamin-borbe.de"))
	list = append(list, http.New("http://slideshow.benjamin-borbe.de").ExpectBody("go.html"))
	list = append(list, http.New("https://www.benjamin-borbe.de/confluence/").ExpectTitle("Dashboard - Confluence"))
	list = append(list, http.New("http://apt.benjamin-borbe.de/bborbe-unstable/Sources").ExpectContent("bborbe-unstable"))
	list = append(list, http.New("http://blog.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	return list
}
