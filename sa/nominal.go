package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	"regexp"
)

const NOMINAL_REGEX = "^[0-9]{1,10}$"

//Nominal is an account nominal code.  Use the constructor
type Nominal string

//NewNominal constructor. Checks that code matches validation pattern
func NewNominal(n string) (Nominal, error) {
	re := regexp.MustCompile(NOMINAL_REGEX)
	if !re.MatchString(n) {
		return "", ErrBadNominal
	}
	return Nominal(n), nil
}

//MustNewNominal constructor. Will Panic if n is not a valid nominal code
func MustNewNominal(n string) Nominal {
	nom, err := NewNominal(n)
	if err != nil {
		panic(err)
	}
	return nom
}

//String implements Stringify interface
func (n Nominal) String() string {
	return string(n)
}
