package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	"time"
)

//SimpleTransactionBuilder a builder for a Simple (2 entry) SplitTransaction
type SimpleTransactionBuilder struct {
	txn *SplitTransaction
}

//NewSimpleTransactionBuilder returns a SimpleTransactionBuilder
func NewSimpleTransactionBuilder(id uint64, drAc, crAc Nominal, amount uint64) *SplitTransactionBuilder {
	ret := &SplitTransactionBuilder{txn: &SplitTransaction{
		txnId:   id,
		date:    time.Now(),
		src:     "",
		ref:     0,
		entries: nil,
		note:    "",
	}}

	return ret.
		WithEntry(*NewEntry(drAc, amount, *NewAcType().Dr())).
		WithEntry(*NewEntry(crAc, amount, *NewAcType().Cr()))
}

//WithDate adds a date to the builder
func (b *SimpleTransactionBuilder) WithDate(dt time.Time) *SimpleTransactionBuilder {
	b.txn.date = dt
	return b
}

//WithNote adds a note to the builder
func (b *SimpleTransactionBuilder) WithNote(note string) *SimpleTransactionBuilder {
	b.txn.note = note
	return b
}

//WithSource adds a source to the builder
func (b *SimpleTransactionBuilder) WithSource(src string) *SimpleTransactionBuilder {
	b.txn.src = src
	return b
}

//WithReference adds a reference to the builder
func (b *SimpleTransactionBuilder) WithReference(ref uint64) *SimpleTransactionBuilder {
	b.txn.ref = ref
	return b
}

//Build builds and returns a SplitTransaction
func (b *SimpleTransactionBuilder) Build() *SplitTransaction {
	return b.txn
}
