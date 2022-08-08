package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import "errors"

var (
	ErrNoChartId             = errors.New("chart id not set")
	ErrBadAccountType        = errors.New("chart contains unknown account type")
	ErrNoChartName           = errors.New("cannot retrieve chart name")
	ErrNoJrnId               = errors.New("no journal id returned")
	ErrBalanceType           = errors.New("cannot determine account type to set balance")
	ErrDummyAccount          = errors.New("no operations available on dummy account type")
	ErrBadXmlParse           = errors.New("error parsing chart definition xml")
	ErrBadNominal            = errors.New("provided value for Nominal does not match pattern: " + NOMINAL_REGEX)
	ErrEntryNotFound         = errors.New("entry not found")
	ErrUnbalancedTransaction = errors.New("transaction is not balanced")
)
