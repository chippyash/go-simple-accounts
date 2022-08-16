//go:build unit

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
)

func TestYouCanGetADebitColumnTitleForAValidAccountType(t *testing.T) {
	at := sa.NewAcType()
	validTypes := map[*sa.AccountType]string{
		at.Dr():        "Debit",
		at.Cr():        "Debit",
		at.Asset():     "Increase",
		at.Bank():      "Increase",
		at.Customer():  "Increase",
		at.Expense():   "Expense",
		at.Income():    "Charge",
		at.Liability(): "Decrease",
		at.Equity():    "Decrease",
		at.Supplier():  "Decrease",
	}
	for ac, expected := range validTypes {
		actual, err := ac.DrTitle()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestYouCanGetACreditColumnTitleForAValidAccountType(t *testing.T) {
	at := sa.NewAcType()
	validTypes := map[*sa.AccountType]string{
		at.Dr():        "Credit",
		at.Cr():        "Credit",
		at.Asset():     "Decrease",
		at.Bank():      "Decrease",
		at.Customer():  "Decrease",
		at.Expense():   "Refund",
		at.Income():    "Income",
		at.Liability(): "Increase",
		at.Equity():    "Increase",
		at.Supplier():  "Increase",
	}
	for ac, expected := range validTypes {
		actual, err := ac.CrTitle()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestYouCanGetBothColumnTitlesForAValidAccountType(t *testing.T) {
	at := sa.NewAcType()
	validTypes := map[*sa.AccountType][]string{
		at.Dr():        {"Debit", "Credit"},
		at.Cr():        {"Debit", "Credit"},
		at.Asset():     {"Increase", "Decrease"},
		at.Bank():      {"Increase", "Decrease"},
		at.Customer():  {"Increase", "Decrease"},
		at.Expense():   {"Expense", "Refund"},
		at.Income():    {"Charge", "Income"},
		at.Liability(): {"Decrease", "Increase"},
		at.Equity():    {"Decrease", "Increase"},
		at.Supplier():  {"Decrease", "Increase"},
	}
	for ac, expected := range validTypes {
		drTitle, crTitle, err := ac.Titles()
		assert.NoError(t, err)
		assert.Equal(t, expected[0], drTitle)
		assert.Equal(t, expected[1], crTitle)
	}
}

func TestFetchingTitlesForInvalidAccountTypeWillReturnError(t *testing.T) {
	//dummy account
	sut := sa.NewAcType()
	_, err := sut.DrTitle()
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrDummyAccount)

	_, err = sut.CrTitle()
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrDummyAccount)

	_, _, err = sut.Titles()
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrDummyAccount)

	//unknown account type
	sut2 := sa.AccountType(0xFF)
	_, err = sut2.DrTitle()
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBalanceType)

	_, err = sut2.CrTitle()
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBalanceType)

	_, _, err = sut2.Titles()
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBalanceType)
}

func TestYouCanGetBalancesForValidAccountTypes(t *testing.T) {
	at := sa.NewAcType()
	validTypes := map[*sa.AccountType][]int64{
		at.Dr():        {2, 1, 1},
		at.Cr():        {1, 2, 1},
		at.Asset():     {2, 1, 1},
		at.Bank():      {2, 1, 1},
		at.Customer():  {2, 1, 1},
		at.Expense():   {2, 1, 1},
		at.Income():    {1, 2, 1},
		at.Liability(): {1, 2, 1},
		at.Equity():    {1, 2, 1},
		at.Supplier():  {1, 2, 1},
	}
	for ac, vals := range validTypes {
		result, err := ac.Balance(vals[0], vals[1])
		assert.NoError(t, err)
		assert.Equal(t, vals[2], result)
	}
}

func TestYouCanGetBalancesForRealAccountTypes(t *testing.T) {
	at := sa.NewAcType().Real()
	realTypes := [][]int64{
		{2, 1, 1},
		{1, 2, 1},
		{2, 2, 0},
	}
	for _, vals := range realTypes {
		result, err := at.Balance(vals[0], vals[1])
		assert.NoError(t, err)
		assert.Equal(t, vals[2], result)
	}
}

func TestBalancesForInvalidAccountTypeWillReturnError(t *testing.T) {
	//dummy account type
	sut := sa.NewAcType()
	_, err := sut.Balance(2, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrDummyAccount)

	//unknown account type
	sut2 := sa.AccountType(0xFF)
	_, err = sut2.Balance(2, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBalanceType)
}
