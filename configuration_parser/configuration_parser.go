package configuration_parser

import (
	"encoding/xml"
	"fmt"

	"time"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_check_http "github.com/bborbe/monitoring/check/http"
	monitoring_check_nop "github.com/bborbe/monitoring/check/nop"
	monitoring_check_tcp "github.com/bborbe/monitoring/check/tcp"
	monitoring_check_webdriver "github.com/bborbe/monitoring/check/webdriver"
	monitoring_node "github.com/bborbe/monitoring/node"
	"github.com/bborbe/webdriver"
)

var logger = log.DefaultLogger

type ConfigurationParser interface {
	ParseConfiguration(content []byte) ([]monitoring_node.Node, error)
}

type configurationParser struct {
	webDriver webdriver.WebDriver
}

type XmlNodes struct {
	NodeList []XmlNode `xml:"node"`
}

type XmlNode struct {
	NodeList         []XmlNode   `xml:"node"`
	ActionList       []XmlAction `xml:"action"`
	Silent           bool        `xml:"silent,attr"`
	Disabled         bool        `xml:"disabled,attr"`
	Check            string      `xml:"check,attr"`
	Port             int         `xml:"port,attr"`
	Retrycount       int         `xml:"retrycount,attr"`
	Timeout          int         `xml:"timeout,attr"`
	Host             string      `xml:"host,attr"`
	Url              string      `xml:"url,attr"`
	Name             string      `xml:"name,attr"`
	ExpectBody       string      `xml:"expectbody,attr"`
	ExpectContent    string      `xml:"expectcontent,attr"`
	ExpectStatusCode int         `xml:"expectstatuscode,attr"`
	ExpectTitle      string      `xml:"expecttitle,attr"`
	Username         string      `xml:"username,attr"`
	Password         string      `xml:"password,attr"`
	PasswordFile     string      `xml:"passwordfile,attr"`
}

type XmlAction struct {
	Type     string        `xml:"type,attr"`
	Value    string        `xml:"value,attr"`
	XPath    string        `xml:"xpath,attr"`
	Duration time.Duration `xml:"duration,attr"`
}

func New(webDriver webdriver.WebDriver) *configurationParser {
	c := new(configurationParser)
	c.webDriver = webDriver
	return c
}

func (c *configurationParser) ParseConfiguration(content []byte) ([]monitoring_node.Node, error) {
	logger.Debugf("parse configuration")
	if len(content) == 0 {
		return nil, fmt.Errorf("can't parse empty content")
	}
	var data XmlNodes
	err := xml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return c.convertXmlNodesToNodes(data.NodeList)
}

func (c *configurationParser) convertXmlNodesToNodes(xmlNodes []XmlNode) ([]monitoring_node.Node, error) {
	var result []monitoring_node.Node
	for _, xmlNode := range xmlNodes {
		node, err := c.convertXmlNodeToNode(xmlNode)
		if err != nil {
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func (c *configurationParser) convertXmlNodeToNode(xmlNode XmlNode) (monitoring_node.Node, error) {
	check, err := c.createCheck(xmlNode)
	if err != nil {
		return nil, err
	}
	nodes, err := c.convertXmlNodesToNodes(xmlNode.NodeList)
	if err != nil {
		return nil, err
	}
	result := monitoring_node.New(check, nodes...).Silent(xmlNode.Silent).Disabled(xmlNode.Disabled)
	return result, nil
}

func (c *configurationParser) createCheck(xmlNode XmlNode) (monitoring_check.Check, error) {
	switch xmlNode.Check {
	case "nop":
		return monitoring_check_nop.New(xmlNode.Name), nil
	case "tcp":
		check := monitoring_check_tcp.New(xmlNode.Host, xmlNode.Port)
		if xmlNode.Timeout > 0 {
			check.Timeout(time.Duration(xmlNode.Timeout) * time.Second)
		}
		if xmlNode.Retrycount > 0 {
			check.RetryCounter(xmlNode.Retrycount)
		}
		return check, nil
	case "http":
		check := monitoring_check_http.New(xmlNode.Url)
		if xmlNode.Timeout > 0 {
			check.Timeout(time.Duration(xmlNode.Timeout) * time.Second)
		}
		if len(xmlNode.ExpectContent) > 0 {
			check.ExpectContent(xmlNode.ExpectContent)
		}
		if len(xmlNode.ExpectBody) > 0 {
			check.ExpectBody(xmlNode.ExpectBody)
		}
		if xmlNode.ExpectStatusCode > 0 {
			check.ExpectStatusCode(xmlNode.ExpectStatusCode)
		}
		if len(xmlNode.ExpectTitle) > 0 {
			check.ExpectTitle(xmlNode.ExpectTitle)
		}
		if len(xmlNode.Username) > 0 && len(xmlNode.Password) > 0 {
			check.Auth(xmlNode.Username, xmlNode.Password)
		}
		if len(xmlNode.Username) > 0 && len(xmlNode.PasswordFile) > 0 {
			check.AuthFile(xmlNode.Username, xmlNode.PasswordFile)
		}
		return check, nil
	case "webdriver":
		check := monitoring_check_webdriver.New(c.webDriver, xmlNode.Url)
		for _, action := range xmlNode.ActionList {
			switch action.Type {
			case "expecttitle":
				check.ExpectTitle(action.Value)
			case "printsource":
				check.PrintSource()
			case "fill":
				check.Fill(action.XPath, action.Value, action.Duration * time.Millisecond)
			case "submit":
				check.Submit(action.XPath, action.Duration * time.Millisecond)
			case "click":
				check.Click(action.XPath, action.Duration * time.Millisecond)
			case "exists":
				check.Exists(action.XPath, action.Duration * time.Millisecond)
			case "notexists":
				check.NotExists(action.XPath, action.Duration * time.Millisecond)
			case "waitfor":
				check.WaitFor(action.XPath, action.Duration * time.Millisecond)
			case "sleep":
				check.Sleep(action.Duration * time.Millisecond)
			default:
				return nil, fmt.Errorf("unkown action '%s'", action.Type)
			}
		}
		return check, nil
	default:
		return nil, fmt.Errorf("not check with typ '%s' found", xmlNode.Check)
	}
}
