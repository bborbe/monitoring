package configuration_parser

import (
	"encoding/xml"
	"fmt"

	"github.com/bborbe/log"
	monitoring_check "github.com/bborbe/monitoring/check"
	monitoring_check_dummy "github.com/bborbe/monitoring/check/dummy"
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
	Silent bool `xml:"silent,attr"`
	Disabled bool `xml:"disabled,attr"`
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
	return convertXmlNodesToNodes(data)
}

func convertXmlNodesToNodes(xmlNodes XmlNodes) ([]monitoring_node.Node, error) {
	var result []monitoring_node.Node
	for _, xmlNode := range xmlNodes.NodeList {
		node, err := convertXmlNodeToNode(xmlNode)
		if err != nil {
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func convertXmlNodeToNode(xmlNode XmlNode) (monitoring_node.Node, error) {
	check := monitoring_check_dummy.New(monitoring_check.NewCheckResultSuccess("ok"), "dummy")
	result := monitoring_node.New(check).Silent(xmlNode.Silent).Disabled(xmlNode.Disabled)
	return result, nil
}
