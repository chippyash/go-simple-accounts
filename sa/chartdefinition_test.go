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

func TestChartDefinition_GetDefinition(t *testing.T) {
	sut := setupChartDefinition(t)
	dom, err := sut.GetDefinition()
	assert.NoError(t, err)
	assert.IsType(t, xmldom.Document{}, *dom)
}

func setupChartDefinition(t *testing.T) *sa.ChartDefinition {
	sutt, err := sa.NewChartDefinition("../tests/_data/personal.xml")
	assert.NoError(t, err)
	return sutt
}
