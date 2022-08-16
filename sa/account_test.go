//go:build unit

package sa_test

import (
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dummy = sa.NewAcType()

var validAccountTypes = []*sa.AccountType{
	dummy.Dr(),
	dummy.Cr(),
	dummy.Asset(),
	dummy.Liability(),
	dummy.Bank(),
	dummy.Customer(),
	dummy.Equity(),
	dummy.Expense(),
	dummy.Income(),
	dummy.Real(),
	dummy.Supplier(),
}

func TestAccount_YouCanCreateAnyValidAccountType(t *testing.T) {
	for _, acc := range validAccountTypes {
		sut := sa.NewAccount("0000", acc, "foo", 0, 0, 1)
		assert.IsType(t, sa.Account{}, *sut)
	}
}

func TestAccount_DummyAccountCantDoAnything(t *testing.T) {
	_, err := dummy.Balance(0, 0)
	assert.Error(t, err)
	_, _, err = dummy.Titles()
	assert.Error(t, err)
	_, err = dummy.DrTitle()
	assert.Error(t, err)
	_, err = dummy.CrTitle()
	assert.Error(t, err)
}

func TestAccount_Name(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	assert.Equal(t, "foo", sut.Name())
}

func TestAccount_Dr(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	assert.Equal(t, int64(0), sut.Dr())
}

func TestAccount_Cr(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	assert.Equal(t, int64(0), sut.Cr())
}

func TestAccount_Balance(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	b, err := sut.Balance()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), b)
}

func TestAccount_Nominal(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	nom := sut.Nominal()
	assert.Equal(t, "0000", nom.String())
}

func TestAccount_Type(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	assert.Equal(t, sa.NewAcType().Asset(), sut.Type())
}

func TestAccount_ChartId(t *testing.T) {
	sut := sa.NewAccount("0000", sa.NewAcType().Asset(), "foo", 0, 0, 1)
	assert.Equal(t, uint64(1), sut.ChartId())
}
