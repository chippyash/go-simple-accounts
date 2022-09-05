//go:build unit
// +build unit

package sa_test

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/stretchr/testify/assert"
	"github.com/subchen/go-xmldom"
	"testing"
)

func TestChartDefinition_NewChartDefinition(t *testing.T) {
	sut, err := sa.NewChartDefinition("../tests/_data/personal.xml")
	assert.NoError(t, err)
	dom, err := sut.GetDefinition()
	assert.NoError(t, err)
	assert.IsType(t, xmldom.Document{}, *dom)
}

func TestChartDefinition_NewChartDefinitionFromString(t *testing.T) {
	def := `<?xml version="1.0" encoding="UTF-8"?>
<chart  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:noNamespaceSchemaLocation="../../sa/chart-definition.xsd"
        name="Personal">
    <account nominal="0000" type="real" name="COA">
        <account nominal="0001" type="dr" name="Balance Sheet"/>
        <account nominal="0002" type="cr" name="Profit And Loss"/>
    </account>
</chart>`

	sut := sa.NewChartDefinitionFromString(def)
	dom, err := sut.GetDefinition()
	assert.NoError(t, err)
	assert.IsType(t, xmldom.Document{}, *dom)
}
