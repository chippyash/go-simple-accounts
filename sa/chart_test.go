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
	"github.com/chippyash/go-hierarchy-tree/tree"
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sut *sa.Chart

func TestChart_GetAccountThatExists(t *testing.T) {
	setupChartTest()
	test := sut.GetAccount(sa.MustNewNominal("1000"))
	assert.NotNil(t, test)
	assert.Equal(t, "Assets", test.Name())
}

func TestChart_GetAccountThatDoesNotExist(t *testing.T) {
	setupChartTest()
	test := sut.GetAccount(sa.MustNewNominal("3000"))
	assert.Nil(t, test)
}

func TestChart_HasAccount(t *testing.T) {
	setupChartTest()
	assert.True(t, sut.HasAccount(sa.MustNewNominal("1000")))
	assert.False(t, sut.HasAccount(sa.MustNewNominal("3000")))
}

func TestChart_GetParentId(t *testing.T) {
	setupChartTest()
	assert.Equal(t, sa.Nominal("0000"), sut.GetParentId(sa.MustNewNominal("1000")))
	assert.Equal(t, sa.Nominal(""), sut.GetParentId(sa.MustNewNominal("3000")))
}

func TestChart_SetRootNode(t *testing.T) {
	setupChartTest()
	newRootNode := tree.NewNode(
		sa.NewAccount(sa.MustNewNominal("9999"), sa.NewAcType().Real(), "Root", 0, 0),
		nil,
	)
	sut.SetRootNode(newRootNode)
	tre := sut.Tree()
	assert.Equal(t, sa.Nominal("9999"), tre.GetValue().(*sa.Account).Nominal())
	assert.Equal(t, 0, len(tre.GetChildren()))
}

func TestChart_Name(t *testing.T) {
	setupChartTest()
	assert.Equal(t, "Test", sut.Name())
}

func TestChart_Crcy(t *testing.T) {
	setupChartTest()
	assert.Equal(t, "GBP", sut.Crcy())
}

func TestChart_Tree(t *testing.T) {
	setupChartTest()
	tre := sut.Tree()
	_, ok := tre.(tree.NodeIFace)
	assert.True(t, ok)
}

func setupChartTest() {
	tr := tree.NewNode(
		sa.NewAccount(sa.MustNewNominal("0000"), sa.NewAcType().Real(), "COA", 0, 0),
		&[]tree.NodeIFace{
			tree.NewNode(
				sa.NewAccount(sa.MustNewNominal("1000"), sa.NewAcType().Asset(), "Assets", 0, 0),
				nil,
			),
			tree.NewNode(
				sa.NewAccount(sa.MustNewNominal("2000"), sa.NewAcType().Liability(), "Liability", 0, 0),
				nil,
			),
		},
	)

	sut = sa.NewChart(1, "Test", "GBP", tr)
}
