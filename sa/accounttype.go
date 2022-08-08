package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	"math"
)

type AccountType uint16

var (
	//Bit values for account types
	dummy     AccountType = 0b0000000000000000 //0 - A dummy account - used internally, not for public consumption
	real      AccountType = 0b0000000000000001 //1 - Base of all accounts - used for root accounts: Balance = abs(cr - dr)
	dr        AccountType = 0b0000000000000011 //3 - Debit account: Balance = dr - cr
	asset     AccountType = 0b0000000000001011 //11 - An account showing assets. Value coming in is DR, going out is CR.
	bank      AccountType = 0b0000000000011011 //27 - An account at a bank.  It is a special form of Asset Account
	customer  AccountType = 0b0000000000101011 //44 - An asset account recording sales to a customer.
	expense   AccountType = 0b0000000001001011 //77 - An account showing destination of expenses.  Expense is shown as DR, refund of expense as CR.
	cr        AccountType = 0b0000000000000101 //5 - Credit account: Balance = cr - dr
	liability AccountType = 0b0000000010000101 //133 - An account recording liabilities (money owing to third parties.) Liability recorded as CR.
	income    AccountType = 0b0000000110000101 //389 - An account showing sources of income.  Income is shown as CR, Refund as DR
	equity    AccountType = 0b0000001010000101 //645 - An account recording the capital or equity of an organisation.  Positive value is shown as CR, negative as DR.  Essentially a form of Liability as it is owed to the shareholders or owners.
	supplier  AccountType = 0b0000010010000101 //1157 - A liability account recording details of purchases from Suppliers.

	//Account display titles for Dr and Cr columns for AccountType
	titles = map[AccountType]map[string]string{
		dr:        {"dr": "Debit", "cr": "Credit"},
		cr:        {"dr": "Debit", "cr": "Credit"},
		asset:     {"dr": "Increase", "cr": "Decrease"},
		bank:      {"dr": "Increase", "cr": "Decrease"},
		customer:  {"dr": "Increase", "cr": "Decrease"},
		expense:   {"dr": "Expense", "cr": "Refund"},
		income:    {"dr": "Charge", "cr": "Income"},
		liability: {"dr": "Decrease", "cr": "Increase"},
		equity:    {"dr": "Decrease", "cr": "Increase"},
		supplier:  {"dr": "Decrease", "cr": "Increase"},
	}

	//map of account type names to value
	named = map[string]*AccountType{
		"REAL":      &real,
		"DR":        &dr,
		"CR":        &cr,
		"ASSET":     &asset,
		"BANK":      &bank,
		"CUSTOMER":  &customer,
		"EXPENSE":   &expense,
		"INCOME":    &income,
		"LIABILITY": &liability,
		"EQUITY":    &equity,
		"SUPPLIER":  &supplier,
	}

	//map of account type values to their names
	values = map[AccountType]string{
		real:      "REAL",
		dr:        "DR",
		cr:        "CR",
		asset:     "ASSET",
		bank:      "BANK",
		customer:  "CUSTOMER",
		expense:   "EXPENSE",
		income:    "INCOME",
		liability: "LIABILITY",
		equity:    "EQUITY",
		supplier:  "SUPPLIER",
	}
)

//NewAcType returns an AccountType for further operations
func NewAcType() *AccountType {
	return &dummy
}

//Real returns a real AccountType
func (a *AccountType) Real() *AccountType {
	return &real
}

//Dr returns a dr AccountType
func (a *AccountType) Dr() *AccountType {
	return &dr
}

//Cr returns a cr AccountType
func (a *AccountType) Cr() *AccountType {
	return &cr
}

//Asset returns an asset AccountType
func (a *AccountType) Asset() *AccountType {
	return &asset
}

//Bank returns a bank AccountType
func (a *AccountType) Bank() *AccountType {
	return &bank
}

//Customer returns a customer AccountType
func (a *AccountType) Customer() *AccountType {
	return &customer
}

//Expense returns an expense AccountType
func (a *AccountType) Expense() *AccountType {
	return &expense
}

//Income returns an income AccountType
func (a *AccountType) Income() *AccountType {
	return &income
}

//Liability returns a liability AccountType
func (a *AccountType) Liability() *AccountType {
	return &liability
}

//Equity returns an equity AccountType
func (a *AccountType) Equity() *AccountType {
	return &equity
}

//Supplier returns a supplier AccountType
func (a *AccountType) Supplier() *AccountType {
	return &supplier
}

//Titles returns the Debit and Credit titles for the AccountType
func (a *AccountType) Titles() (string, string, error) {
	if *a == dummy {
		return "", "", ErrDummyAccount
	}
	if _, ok := titles[*a]; !ok {
		return "", "", ErrBalanceType
	}
	return titles[*a]["dr"], titles[*a]["cr"], nil
}

//DrTitle returns the Debit title for the account type
func (a *AccountType) DrTitle() (string, error) {
	if *a == dummy {
		return "", ErrDummyAccount
	}
	v, _, err := a.Titles()
	return v, err
}

//CrTitle returns the Credit title for the account type
func (a *AccountType) CrTitle() (string, error) {
	if *a == dummy {
		return "", ErrDummyAccount
	}
	_, v, err := a.Titles()
	return v, err

}

//Balance returns the balance of two values dependent on the account type
func (a *AccountType) Balance(drVal, crVal uint64) (uint64, error) {
	if *a == dummy {
		return 0, ErrDummyAccount
	}
	if _, ok := titles[*a]; !ok {
		return 0, ErrBalanceType
	}
	if *a&dr == dr {
		//debit account type
		return drVal - crVal, nil
	}
	if *a&cr == cr {
		//credit account type
		return crVal - drVal, nil
	}
	if *a&real == real {
		//real balance - should always be zero as it is the root account
		return uint64(math.Abs(float64(drVal) - float64(crVal))), nil
	}
	return 0, ErrBalanceType
}

func GetNamedAccountTypes() map[string]*AccountType {
	return named
}

func GetValuedAccountTypes() map[AccountType]string {
	return values
}
