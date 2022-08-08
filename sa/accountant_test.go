//go:build integration
// +build integration

package sa_test

import (
	"database/sql"
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var accountant *sa.Accountant
var db *sql.DB

func TestAccountant_CreateChart(t *testing.T) {
	setupAccountantTest(t)
	def, err := sa.NewChartDefinition("../tests/_data/personal.xml")
	assert.NoError(t, err)
	lastId, err := accountant.CreateChart("Test", "GBP", def)
	assert.NoError(t, err)
	assert.True(t, lastId > 0)
	teardownAccountantTest(t)
}

func TestAccountant_FetchChart(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	lastId, _ := accountant.CreateChart("Test", "GBP", def)
	accountant = sa.NewAccountant(db, lastId, "GBP")
	chart, err := accountant.FetchChart()
	assert.NoError(t, err)
	assert.Equal(t, 5, chart.Tree().GetHeight())
	teardownAccountantTest(t)
}

func TestAccountant_WriteTransactionWithDate(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	_, _ = accountant.CreateChart("Test", "GBP", def)
	dt, err := time.Parse(time.RFC3339, "2020-08-05T14:36:00+01:00")
	assert.NoError(t, err)
	txn := sa.NewSimpleTransactionBuilder(0, "0001", "0002", 100).
		WithDate(dt).
		Build()
	jrnId, err := accountant.WriteTransactionWithDate(txn, dt)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), jrnId)
	teardownAccountantTest(t)
}

func TestAccountant_FetchTransaction(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	_, _ = accountant.CreateChart("Test", "GBP", def)
	dt, _ := time.Parse(time.RFC3339, "2020-08-05T14:36:00+01:00")
	txn := sa.NewSimpleTransactionBuilder(0, "0001", "0002", 100).
		WithDate(dt).
		Build()
	jrnId, _ := accountant.WriteTransactionWithDate(txn, dt)

	journal, err := accountant.FetchTransaction(jrnId)
	assert.NoError(t, err)
	//NB above date was set as BST - the stored date is UTC - that is correct
	assert.Equal(t, "2020-08-05T13:36:00Z", journal.Date().Format(time.RFC3339))

	teardownAccountantTest(t)
}

func TestAccountant_FetchAccountJournals(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	_, _ = accountant.CreateChart("Test", "GBP", def)
	dt, _ := time.Parse(time.RFC3339, "2020-08-05T14:36:00+01:00")
	txn := sa.NewSimpleTransactionBuilder(0, "0001", "0002", 100).
		WithDate(dt).
		Build()
	_, _ = accountant.WriteTransactionWithDate(txn, dt)
	txn = sa.NewSimpleTransactionBuilder(0, "0002", "0001", 10).
		WithDate(dt).
		Build()
	_, _ = accountant.WriteTransactionWithDate(txn, dt)

	entries, err := accountant.FetchAccountJournals("0001")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entries))
	assert.Equal(t, 1, len(entries[0].Entries()))
	assert.Equal(t, 1, len(entries[1].Entries()))

	teardownAccountantTest(t)
}

func TestAccountant_AddAccount(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	_, _ = accountant.CreateChart("Test", "GBP", def)

	//2nd root account - error
	nom, _ := sa.NewNominal("1111")
	err := accountant.AddAccount(nom, sa.NewAcType().Asset(), "foo", nil)
	assert.Error(t, err)

	//parent doesn't exist - error
	prnt, _ := sa.NewNominal("9999")
	err = accountant.AddAccount(nom, sa.NewAcType().Asset(), "foo", &prnt)
	assert.Error(t, err)

	//valid insertion
	prnt, _ = sa.NewNominal("0000")
	err = accountant.AddAccount(nom, sa.NewAcType().Asset(), "foo", &prnt)
	assert.NoError(t, err)

	teardownAccountantTest(t)
}

func TestAccountant_DelAccount(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	_, _ = accountant.CreateChart("Test", "GBP", def)

	txn := sa.NewSimpleTransactionBuilder(0, "0001", "0002", 100).
		Build()
	_, _ = accountant.WriteTransaction(txn)

	//attempt to delete account with non-zero balance - error
	nom, _ := sa.NewNominal("0001")
	err := accountant.DelAccount(nom)
	assert.Error(t, err)

	//zero balance - ok
	nom, _ = sa.NewNominal("1210")
	err = accountant.DelAccount(nom)
	assert.NoError(t, err)

	teardownAccountantTest(t)
}

func TestAccountant_NextNominal(t *testing.T) {
	setupAccountantTest(t)
	def, _ := sa.NewChartDefinition("../tests/_data/personal.xml")
	_, _ = accountant.CreateChart("Test", "GBP", def)

	//ac 1800 - Equipment has no children
	starter, _ := sa.NewNominal("1801")
	prnt, _ := sa.NewNominal("1800")

	next, err := accountant.NextNominal(prnt, starter)
	assert.NoError(t, err)
	assert.Equal(t, "1801", next.String())

	//add a child account
	err = accountant.AddAccount(starter, sa.NewAcType().Asset(), "Tractor", &prnt)
	assert.NoError(t, err)

	next, err = accountant.NextNominal(prnt, starter)
	assert.NoError(t, err)
	assert.Equal(t, "1802", next.String())

	teardownAccountantTest(t)
}

func setupAccountantTest(t *testing.T) {
	config := mysql.Config{
		User:                 os.Getenv("DBUID"),
		Passwd:               os.Getenv("DBPWD"),
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	dba, err := sql.Open("mysql", config.FormatDSN())
	assert.NoError(t, err)
	accountant = sa.NewAccountant(dba, 0, "GBP")
	db = dba
}

func teardownAccountantTest(t *testing.T) {
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
