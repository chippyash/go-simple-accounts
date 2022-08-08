package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

type Entry struct {
	entryId *Nominal
	amount  uint64
	tpe     *AccountType
}

//NewEntry constructor for Entry
func NewEntry(entryId Nominal, amount uint64, tpe AccountType) *Entry {
	return &Entry{
		entryId: &entryId,
		amount:  amount,
		tpe:     &tpe,
	}
}

//Id returns the Entry Id
func (e *Entry) Id() *Nominal {
	return e.entryId
}

//Type returns the Entry Type
func (e *Entry) Type() *AccountType {
	return e.tpe
}

//Amount returns the Entry Amount
func (e *Entry) Amount() uint64 {
	return e.amount
}
