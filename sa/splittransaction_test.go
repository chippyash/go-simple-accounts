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

func TestNewSplitTransactionBuilder_Defaults(t *testing.T) {
	sut := sa.NewSplitTransactionBuilder(0).Build()
	assert.Equal(t, uint64(0), sut.Id())
	assert.IsType(t, time.Time{}, sut.Date())
	assert.Equal(t, "", sut.Src())
	assert.Equal(t, "", sut.Note())
	assert.Equal(t, uint64(0), sut.Ref())
	amt, err := sut.GetAmount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), amt)
	assert.True(t, sut.CheckBalance())
}

func TestSplitTransactionBuilder_WithDate(t *testing.T) {
	dt, err := time.Parse(time.RFC3339, "2020-08-05T14:36:00+01:00")
	assert.NoError(t, err)
	sut := sa.NewSplitTransactionBuilder(0).WithDate(dt).Build()
	assert.Equal(t, dt, sut.Date())
}

func TestSplitTransactionBuilder_WithNote(t *testing.T) {
	sut := sa.NewSplitTransactionBuilder(0).WithNote("foo").Build()
	assert.Equal(t, "foo", sut.Note())
}

func TestSplitTransactionBuilder_WithReference(t *testing.T) {
	sut := sa.NewSplitTransactionBuilder(0).WithReference(1).Build()
	assert.Equal(t, uint64(1), sut.Ref())
}

func TestSplitTransactionBuilder_WithSource(t *testing.T) {
	sut := sa.NewSplitTransactionBuilder(0).WithSource("src").Build()
	assert.Equal(t, "src", sut.Src())
}

func TestSplitTransactionBuilder_WithEntry(t *testing.T) {
	nom := sa.MustNewNominal("1000")
	sut := sa.NewSplitTransactionBuilder(0).WithEntry(*sa.NewEntry(nom, 100, *sa.NewAcType().Dr())).Build()
	entries := sut.Entries()
	assert.Equal(t, 1, len(entries))
	assert.False(t, sut.CheckBalance())
}

func TestSplitTransactionBuilder_WithEntries(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entry2 := sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())
	entries := sa.Entries{entry1, entry2}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	entries = sut.Entries()
	assert.Equal(t, 2, len(entries))
	assert.True(t, sut.CheckBalance())
}

func TestSplitTransaction_GetAmountBalancedTransaction(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entry2 := sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())
	entries := sa.Entries{entry1, entry2}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	amt, err := sut.GetAmount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), amt)
}

func TestSplitTransaction_GetAmountUnbalancedTransaction(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entries := sa.Entries{entry1}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	_, err := sut.GetAmount()
	assert.Error(t, err)
}

func TestSplitTransaction_GetDrAc(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entry2 := sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())
	entries := sa.Entries{entry1, entry2}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	noms := sut.GetDrAc()
	assert.Equal(t, 1, len(noms))
	assert.Equal(t, "1000", noms[0].String())
}

func TestSplitTransaction_GetCrAc(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entry2 := sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())
	entries := sa.Entries{entry1, entry2}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	noms := sut.GetCrAc()
	assert.Equal(t, 1, len(noms))
	assert.Equal(t, "2000", noms[0].String())
}

func TestSplitTransaction_IsSimple(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entry2 := sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())
	entries := sa.Entries{entry1, entry2}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	assert.True(t, sut.IsSimple())
}

func TestSplitTransaction_IsNotSimple(t *testing.T) {
	entry1 := sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())
	entry2 := sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())
	entry3 := sa.NewEntry(sa.MustNewNominal("3000"), 100, *sa.NewAcType().Dr())
	entry4 := sa.NewEntry(sa.MustNewNominal("4000"), 100, *sa.NewAcType().Cr())
	entries := sa.Entries{entry1, entry2, entry3, entry4}
	sut := sa.NewSplitTransactionBuilder(0).WithEntries(entries).Build()
	assert.False(t, sut.IsSimple())
}
