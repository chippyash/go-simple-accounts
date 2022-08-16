//go:build integration

package tests_test

/**
Test capabilities of managing multiple charts of account in same database
*/

import (
	"database/sql"
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var accountant1 *sa.Accountant
var accountant2 *sa.Accountant
var db *sql.DB

func TestYouCanCreateMultipleCharts(t *testing.T) {
	setupMultiChartTest(t)
	chart1, err := accountant1.FetchChart()
	assert.NoError(t, err)
	chart2, err := accountant2.FetchChart()
	assert.NoError(t, err)
	assert.NotEqual(t, chart1.Id(), chart2.Id())
	assert.NotEqual(t, chart1.Name(), chart2.Name())

	teardownMultiChartTestTest(t)
}

func TestSettingTransactionsForMultiChartsDoNotOverlap(t *testing.T) {
	setupMultiChartTest(t)

	dt, _ := time.Parse(time.RFC3339, "2020-08-05T14:36:00+01:00")
	txn := sa.NewSimpleTransactionBuilder(0, "1210", "4100", 100).
		WithDate(dt).
		Build()
	_, _ = accountant1.WriteTransactionWithDate(txn, dt)
	txn = sa.NewSimpleTransactionBuilder(0, "6120", "1210", 10).
		WithDate(dt).
		Build()
	_, _ = accountant1.WriteTransactionWithDate(txn, dt)

	txn = sa.NewSimpleTransactionBuilder(0, "1210", "4100", 200).
		WithDate(dt).
		Build()
	_, _ = accountant2.WriteTransactionWithDate(txn, dt)
	txn = sa.NewSimpleTransactionBuilder(0, "6120", "1210", 20).
		WithDate(dt).
		Build()
	_, _ = accountant2.WriteTransactionWithDate(txn, dt)

	values1 := map[string][]int64{
		"0000": {110, 110},
		"0001": {100, 10},
		"1000": {100, 10},
		"1100": {100, 10},
		"1200": {100, 10},
		"1210": {100, 10},
		"0002": {10, 100},
		"4000": {0, 100},
		"4100": {0, 100},
		"6000": {10, 0},
		"6100": {10, 0},
		"6120": {10, 0},
	}
	chart1, _ := accountant1.FetchChart()
	for nom, vals := range values1 {
		ac := chart1.GetAccount(sa.MustNewNominal(nom))
		assert.Equal(t, vals[0], ac.Dr(), "DR a/c value: %d not equal %d for Nominal: %s", ac.Dr(), vals[0], nom)
		assert.Equal(t, vals[1], ac.Cr(), "CR a/c value: %d not equal %d for Nominal: %s", ac.Cr(), vals[1], nom)
	}

	values2 := map[string][]int64{
		"0000": {220, 220},
		"0001": {200, 20},
		"1000": {200, 20},
		"1100": {200, 20},
		"1200": {200, 20},
		"1210": {200, 20},
		"0002": {20, 200},
		"4000": {0, 200},
		"4100": {0, 200},
		"6000": {20, 0},
		"6100": {20, 0},
		"6120": {20, 0},
	}
	chart2, _ := accountant2.FetchChart()
	for nom, vals := range values2 {
		ac := chart2.GetAccount(sa.MustNewNominal(nom))
		assert.Equal(t, vals[0], ac.Dr(), "DR a/c value: %d not equal %d for Nominal: %s", ac.Dr(), vals[0], nom)
		assert.Equal(t, vals[1], ac.Cr(), "CR a/c value: %d not equal %d for Nominal: %s", ac.Cr(), vals[1], nom)
	}

	teardownMultiChartTestTest(t)
}

func setupMultiChartTest(t *testing.T) {
	config := mysql.Config{
		User:                 os.Getenv("DBUID"),
		Passwd:               os.Getenv("DBPWD"),
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	dba, err := sql.Open("mysql", config.FormatDSN())
	assert.NoError(t, err)
	accountant1 = sa.NewAccountant(dba, 0, "GBP")
	accountant2 = sa.NewAccountant(dba, 0, "GBP")

	def1, err := sa.NewChartDefinition("../tests/_data/personal.xml")
	assert.NoError(t, err)
	def2, err := sa.NewChartDefinition("../tests/_data/personal.xml")
	assert.NoError(t, err)
	lastId, err := accountant1.CreateChart("Test 1", "GBP", def1)
	assert.NoError(t, err)
	assert.True(t, lastId > 0)
	lastId2, err := accountant2.CreateChart("Test 2", "GBP", def2)
	assert.NoError(t, err)
	assert.True(t, lastId2 > 0)
	assert.True(t, lastId != lastId2)

	db = dba
}

func teardownMultiChartTestTest(t *testing.T) {
	_, err := db.Exec("delete from sa_coa")
	assert.NoError(t, err)
	_, err = db.Exec("alter table sa_coa AUTO_INCREMENT=1")
	assert.NoError(t, err)
	_, err = db.Exec("delete from sa_journal")
	assert.NoError(t, err)
	_, err = db.Exec("alter table sa_journal AUTO_INCREMENT=1")
	assert.NoError(t, err)
	_, err = db.Exec("alter table sa_journal_entry AUTO_INCREMENT=1")
	assert.NoError(t, err)
	_, err = db.Exec("delete from sa_coa_ledger")
	assert.NoError(t, err)
	_, err = db.Exec("alter table sa_coa_ledger AUTO_INCREMENT=1")
	assert.NoError(t, err)
}
