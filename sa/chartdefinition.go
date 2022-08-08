package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	_ "embed"
	"github.com/jbussdieker/golibxml"
	"github.com/krolaw/xsd"
	"github.com/subchen/go-xmldom"
	"os"
	"unsafe"
)

//go:embed chart-definition.xsd
var xsdSchema []byte

//ChartDefinition is a helper to retrieve chart definition xml
type ChartDefinition struct {
	xmlFileName string
}

//NewChartDefinition constructor
func NewChartDefinition(xmlFileName string) (*ChartDefinition, error) {
	_, err := os.Stat(xmlFileName)
	if err != nil {
		return nil, err
	}
	return &ChartDefinition{xmlFileName: xmlFileName}, nil
}

//GetDefinition returns parsed xml as Dom Document
func (c *ChartDefinition) GetDefinition() (*xmldom.Document, error) {
	//_, err := c.validate()
	//if err != nil {
	//	return nil, err
	//}

	doc, err := xmldom.ParseFile(c.xmlFileName)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (c *ChartDefinition) validate() (bool, error) {
	schema, err := xsd.ParseSchema(xsdSchema)
	if err != nil {
		return false, err
	}

	doc := golibxml.ParseFile(c.xmlFileName)
	if doc == nil {
		return false, ErrBadXmlParse
	}

	if err = schema.Validate(xsd.DocPtr((unsafe.Pointer(doc.Ptr)))); err != nil {
		return false, err
	}

	return true, nil
}
