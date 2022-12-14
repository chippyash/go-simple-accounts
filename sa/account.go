package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

//Account is an account in a Chart
type Account struct {
	nominal Nominal
	tpe     *AccountType
	name    string
	acDr    int64
	acCr    int64
	chartId uint64
}

//NewAccount Account constructor
func NewAccount(nominal Nominal, tpe *AccountType, name string, acDr, acCr int64, chartId uint64) *Account {
	return &Account{
		nominal: nominal,
		tpe:     tpe,
		name:    name,
		acDr:    acDr,
		acCr:    acCr,
		chartId: chartId,
	}
}

//Dr retrieve account debit amount
func (a *Account) Dr() int64 {
	return a.acDr
}

//Cr retrieve account credit amount
func (a *Account) Cr() int64 {
	return a.acCr
}

//Balance retrieves account balance
func (a *Account) Balance() (int64, error) {
	return a.tpe.Balance(a.acDr, a.acCr)
}

//Nominal retrieve account nominal code (id)
func (a *Account) Nominal() Nominal {
	return a.nominal
}

//Name retrieve account name
func (a *Account) Name() string {
	return a.name
}

//Type retrieve account type
func (a *Account) Type() *AccountType {
	return a.tpe
}

func (a *Account) ChartId() uint64 {
	return a.chartId
}
