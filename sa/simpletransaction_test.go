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
	"testing"
	"time"
)

func TestNewSimpleTransactionBuilder_Defaults(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).Build()
	assert.Equal(t, uint64(0), sut.Id())
	assert.IsType(t, time.Time{}, sut.Date())
	assert.Equal(t, "", sut.Src())
	assert.Equal(t, "", sut.Note())
	assert.Equal(t, uint64(0), sut.Ref())
	amt, err := sut.GetAmount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), amt)
	assert.True(t, sut.CheckBalance())
}

func TestSimpleTransactionBuilder_WithDate(t *testing.T) {
	dt, err := time.Parse(time.RFC3339, "2020-08-05T14:36:00+01:00")
	assert.NoError(t, err)
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		WithDate(dt).
		Build()
	assert.Equal(t, dt, sut.Date())
}

func TestSimpleTransactionBuilder_WithNote(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		WithNote("foo").
		Build()
	assert.Equal(t, "foo", sut.Note())
}

func TestSimpleTransactionBuilder_WithReference(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		WithReference(1).
		Build()
	assert.Equal(t, uint64(1), sut.Ref())
}

func TestSimpleTransactionBuilder_WithSource(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		WithSource("src").
		Build()
	assert.Equal(t, "src", sut.Src())
}

func TestSimpleTransaction_GetAmount(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		Build()
	amt, err := sut.GetAmount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), amt)
}

func TestSimpleTransaction_GetDrAc(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		Build()
	noms := sut.GetDrAc()
	assert.Equal(t, 1, len(noms))
	assert.Equal(t, "1000", noms[0].String())
}

func TestSimpleTransaction_GetCrAc(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		Build()
	noms := sut.GetCrAc()
	assert.Equal(t, 1, len(noms))
	assert.Equal(t, "2000", noms[0].String())
}

func TestSimpleTransaction_IsSimple(t *testing.T) {
	sut := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
		Build()
	assert.True(t, sut.IsSimple())
}
