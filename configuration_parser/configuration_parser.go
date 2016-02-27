package configuration_parser

import (
	"encoding/xml"
	"fmt"

	"time"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_check_http "github.com/bborbe/monitoring/check/http"
	monitoring_check_tcp "github.com/bborbe/monitoring/check/tcp"
	monitoring_check_webdriver "github.com/bborbe/monitoring/check/webdriver"
	monitoring_node "github.com/bborbe/monitoring/node"
)

var logger = log.DefaultLogger

type ConfigurationParser interface {
	ParseConfiguration(content []byte) ([]monitoring_node.Node, error)
}

type configurationParser struct {
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
	ExpectBody       string      `xml:"expectbody,attr"`
	ExpectContent    string      `xml:"expectcontent,attr"`
	ExpectStatusCode int         `xml:"expectstatuscode,attr"`
	ExpectTitle      string      `xml:"expecttitle,attr"`
	Username         string      `xml:"username,attr"`
	Password         string      `xml:"password,attr"`
	PasswordFile     string      `xml:"passwordfile,attr"`
}

type XmlAction struct {
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
	XPath string `xml:"xpath,attr"`
}

func New() *configurationParser {
	return new(configurationParser)
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
	return convertXmlNodesToNodes(data.NodeList)
}

func convertXmlNodesToNodes(xmlNodes []XmlNode) ([]monitoring_node.Node, error) {
	var result []monitoring_node.Node
	for _, xmlNode := range xmlNodes {
		node, err := convertXmlNodeToNode(xmlNode)
		if err != nil {
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func convertXmlNodeToNode(xmlNode XmlNode) (monitoring_node.Node, error) {
	check, err := createCheck(xmlNode)
	if err != nil {
		return nil, err
	}
	nodes, err := convertXmlNodesToNodes(xmlNode.NodeList)
	if err != nil {
		return nil, err
	}
	result := monitoring_node.New(check, nodes...).Silent(xmlNode.Silent).Disabled(xmlNode.Disabled)
	return result, nil
}

func createCheck(xmlNode XmlNode) (monitoring_check.Check, error) {
	switch xmlNode.Check {
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
		check := monitoring_check_webdriver.New(xmlNode.Url)
		for _, action := range xmlNode.ActionList {
			switch action.Type {
			case "expecttitle":
				check.ExpectTitle(action.Value)
			case "fill":
				check.Fill(action.XPath, action.Value)
			case "submit":
				check.Submit(action.XPath)
			default:
				return nil, fmt.Errorf("unkown action '%s'", action.Type)
			}
		}
		return check, nil
	default:
		return nil, fmt.Errorf("not check with typ '%s' found", xmlNode.Check)
	}
}
