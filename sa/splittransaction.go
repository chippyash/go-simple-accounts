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

//SplitTransaction is a basic multi ledger journal entry
type SplitTransaction struct {
	txnId   uint64
	date    time.Time
	note    string
	src     string
	ref     uint64
	entries Entries
}

//Id returns the transaction id
func (s *SplitTransaction) Id() uint64 {
	return s.txnId
}

//Date returns the transaction datetime
func (s *SplitTransaction) Date() time.Time {
	return s.date
}

//Note returns the transaction note
func (s *SplitTransaction) Note() string {
	return s.note
}

//Src returns the transaction src reference
func (s *SplitTransaction) Src() string {
	return s.src
}

//Ref returns the transaction reference id
func (s *SplitTransaction) Ref() uint64 {
	return s.ref
}

//Entries returns the transaction ledger journal entries
func (s *SplitTransaction) Entries() Entries {
	return s.entries
}

//CheckBalance returns true if the transaction entries balance else false
func (s *SplitTransaction) CheckBalance() bool {
	return s.entries.CheckBalance()
}

//GetEntry returns the entry for the nominal code
func (s *SplitTransaction) GetEntry(nominal Nominal) (*Entry, error) {
	ret := s.entries.Filter(func(entry *Entry) bool {
		return *entry.Id() == nominal
	})
	if len(ret) != 1 {
		return nil, ErrEntryNotFound
	}
	return ret[0], nil
}

func (s *SplitTransaction) GetAmount() (int64, error) {
	if !s.CheckBalance() {
		return 0, ErrUnbalancedTransaction
	}
	var tot int64 = 0
	for _, entry := range s.entries {
		tot += entry.Amount()
	}
	return tot / 2, nil
}

//GetDrAc return debit account ids
func (s *SplitTransaction) GetDrAc() []*Nominal {
	drAc := NewAcType().Dr()
	filtered := s.entries.Filter(func(entry *Entry) bool {
		tp := entry.Type()
		return *tp&*drAc == *drAc
	})

	ret := make([]*Nominal, len(filtered))
	for k, v := range filtered {
		ret[k] = v.Id()
	}
	return ret
}

//GetCrAc return credit account ids
func (s *SplitTransaction) GetCrAc() []*Nominal {
	crAc := NewAcType().Cr()
	filtered := s.entries.Filter(func(entry *Entry) bool {
		tp := entry.Type()
		return *tp&*crAc == *crAc
	})

	ret := make([]*Nominal, len(filtered))
	for k, v := range filtered {
		ret[k] = v.Id()
	}
	return ret
}

//IsSimple is this a simple transaction?
func (s *SplitTransaction) IsSimple() bool {
	return len(s.GetDrAc()) == 1 && len(s.GetCrAc()) == 1
}

//SplitTransactionBuilder a build for a Split Transaction
type SplitTransactionBuilder struct {
	txn *SplitTransaction
}

//NewSplitTransactionBuilder returns a SplitTransactionBuilder
func NewSplitTransactionBuilder(id uint64) *SplitTransactionBuilder {
	return &SplitTransactionBuilder{txn: &SplitTransaction{
		txnId:   id,
		date:    time.Now(),
		note:    "",
		src:     "",
		ref:     0,
		entries: nil,
	}}
}

//WithDate adds a date to the builder
func (b *SplitTransactionBuilder) WithDate(dt time.Time) *SplitTransactionBuilder {
	b.txn.date = dt
	return b
}

//WithNote adds a note to the builder
func (b *SplitTransactionBuilder) WithNote(note string) *SplitTransactionBuilder {
	b.txn.note = note
	return b
}

//WithSource adds a source to the builder
func (b *SplitTransactionBuilder) WithSource(src string) *SplitTransactionBuilder {
	b.txn.src = src
	return b
}

//WithReference adds a reference to the builder
func (b *SplitTransactionBuilder) WithReference(ref uint64) *SplitTransactionBuilder {
	b.txn.ref = ref
	return b
}

//WithEntries adds a set of Entry to the builder
func (b *SplitTransactionBuilder) WithEntries(entries Entries) *SplitTransactionBuilder {
	b.txn.entries = append(b.txn.entries, entries...)
	return b
}

//WithEntry adds a single Entry to the builder
func (b *SplitTransactionBuilder) WithEntry(entry Entry) *SplitTransactionBuilder {
	b.txn.entries = append(b.txn.entries, &entry)
	return b
}

//Build builds and returns a SplitTransaction
func (b *SplitTransactionBuilder) Build() *SplitTransaction {
	return b.txn
}
