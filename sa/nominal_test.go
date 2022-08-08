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
	"fmt"
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNominal(t *testing.T) {
	sut, err := sa.NewNominal("0000")
	assert.NoError(t, err)
	assert.Equal(t, sa.Nominal("0000"), sut)
}

func TestNewNominal_OutOfBounds(t *testing.T) {
	//not enough digits
	_, err := sa.NewNominal("")
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBadNominal)
	//too many digits
	_, err = sa.NewNominal("00000000000")
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBadNominal)
	//contains non digits
	_, err = sa.NewNominal("abc")
	assert.Error(t, err)
	assert.ErrorIs(t, err, sa.ErrBadNominal)
}

func TestNominal_Stringify(t *testing.T) {
	sut, _ := sa.NewNominal("0000")
	test := fmt.Sprintf("%s", sut)
	assert.Equal(t, "0000", test)
}
