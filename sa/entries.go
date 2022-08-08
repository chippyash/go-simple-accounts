package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

//Entries is a set of transaction entries
type Entries []*Entry

//CheckBalance returns true if the set of entries balance else false
func (e *Entries) CheckBalance() bool {
	var balance uint64 = 0
	drAc := *NewAcType().Dr()
	for _, entry := range *e {
		if *entry.Type()&drAc == drAc {
			balance -= entry.Amount()
		} else {
			balance += entry.Amount()
		}
	}

	return balance == 0
}

//FilterFunc is a definition for a filter function
type FilterFunc func(entry *Entry) bool

//Filter returns entries with a filter function
func (e *Entries) Filter(f FilterFunc) Entries {
	ret := make(Entries, 0)
	for _, entry := range *e {
		if f(entry) {
			ret = append(ret, entry)
		}
	}

	return ret
}
