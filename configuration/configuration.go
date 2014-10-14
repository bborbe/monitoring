package configuration

import (
	"github.com/bborbe/monitoring/check"
	"github.com/bborbe/monitoring/check/http"
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
	list = append(list, http.New("http://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	list = append(list, http.New("https://www.benjamin-borbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	list = append(list, http.New("http://www.benjaminborbe.de").ExpectTitle("Benjamin Borbe Fotografie"))
	list = append(list, http.New("http://jenkins.benjamin-borbe.de").ExpectTitle("Dashboard [Jenkins]"))
	list = append(list, http.New("http://www.harteslicht.de").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht."))
	list = append(list, http.New("http://www.harteslicht.com").ExpectTitle("www.Harteslicht.com | Fotografieren das Spass macht."))
	list = append(list, http.New("http://kickstart.benjamin-borbe.de").ExpectContent("ks.cfg"))
	list = append(list, http.New("http://ip.benjamin-borbe.de"))
	list = append(list, http.New("http://slideshow.benjamin-borbe.de").ExpectContent("go.html"))
	list = append(list, http.New("https://www.benjamin-borbe.de/confluence/").ExpectTitle("Dashboard - Confluence"))
	return list
}
