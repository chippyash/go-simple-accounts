package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/chippyash/go-hierarchy-tree/tree"
	_ "github.com/go-sql-driver/mysql"
	"github.com/subchen/go-xmldom"
	"strconv"
	"strings"
	"time"
)

//Accountant The main API interface to Simple Accounts
type Accountant struct {
	db      *sql.DB
	chartId uint64
	crcy    string
}

type ledger struct {
	PrntId  uint64
	Id      uint64
	Nominal Nominal
	Name    string
	Tpe     string
	AcDr    uint64
	AcCr    uint64
}
type ledgerLines []ledger

//NewAccountant returns a new Accountant
func NewAccountant(adapter *sql.DB, chartId uint64, crcy string) *Accountant {
	return &Accountant{
		db:      adapter,
		chartId: chartId,
		crcy:    crcy,
	}
}

//CreateChart creates a new chart of accounts from a COA definition file
func (a *Accountant) CreateChart(chartName, crcy string, def *ChartDefinition) (uint64, error) {
	dom, err := def.GetDefinition()
	if err != nil {
		return 0, err
	}
	//create chart tree
	root := dom.Root.Query("/account")[0]
	treeRoot := tree.NewNode(nil, nil)
	err = buildTreeFromXml(treeRoot, root)
	if err != nil {
		return 0, err
	}

	chart := NewChart(0, chartName, crcy, treeRoot)
	chartId, err := a.storeChart(chart)
	if err != nil {
		return 0, err
	}

	errV := treeRoot.Accept(NewNodeSaver(a.db, chartId))
	if errV != nil {
		a.chartId = chartId
		return chartId, errV.(error)
	}

	a.chartId = chartId
	return chartId, nil
}

func buildTreeFromXml(tre tree.NodeIFace, node *xmldom.Node) error {
	//set value of current node
	nom, err := NewNominal(node.GetAttributeValue("nominal"))
	if err != nil {
		return err
	}
	acType, ok := GetNamedAccountTypes()[strings.ToUpper(node.GetAttributeValue("type"))]
	if !ok {
		return ErrBadAccountType
	}
	ac := NewAccount(
		nom,
		acType,
		node.GetAttributeValue("name"),
		0,
		0,
	)
	tre.SetValue(ac)

	//recurse through child accounts
	for _, child := range node.GetChildren("account") {
		childTree := tree.NewNode(nil, nil)
		tre.AddChild(childTree)
		err := buildTreeFromXml(childTree, child)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Accountant) storeChart(chart *Chart) (uint64, error) {
	res, err := a.db.Query("select sa_fu_add_chart(?) as lastId", chart.Name())
	if err != nil {
		return 0, err
	}
	if res.Err() != nil {
		return 0, res.Err()
	}
	var lastId int64 = 0
	defer res.Close()
	if res.Next() {
		err = res.Scan(&lastId)
		return uint64(lastId), err
	}

	return 0, nil
}

//FetchChart fetches a chart from storage
func (a *Accountant) FetchChart() (*Chart, error) {
	if a.chartId == 0 {
		return nil, ErrNoChartId
	}

	//retrieve the ledgers
	res, err := a.db.Query("call sa_sp_get_tree(?)", a.chartId)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	ledgers := make(ledgerLines, 0)
	defer res.Close()
	for res.Next() {
		l := ledger{}
		err = res.Scan(&l.PrntId, &l.Id, &l.Nominal, &l.Name, &l.Tpe, &l.AcDr, &l.AcCr)
		if err != nil {
			return nil, err
		}
		ledgers = append(ledgers, l)
	}

	//retrieve the chart name
	res2, err := a.db.Query("select name from sa_coa where id = ?", a.chartId)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	var chartName string
	defer res2.Close()
	if !res2.Next() {
		return nil, ErrNoChartName
	}
	err = res2.Scan(&chartName)
	if err != nil {
		return nil, err
	}

	//get the root tree node
	rootLedger := ledgers[0]
	ledgers = ledgers[1:]
	acType, ok := GetNamedAccountTypes()[rootLedger.Tpe]
	if !ok {
		return nil, ErrBadAccountType
	}
	root := tree.NewNode(
		NewAccount(
			rootLedger.Nominal,
			acType,
			rootLedger.Name,
			rootLedger.AcDr,
			rootLedger.AcCr,
		),
		nil,
	)

	root, err = buildTreeFromDb(root, ledgers, rootLedger.Id)
	if err != nil {
		return nil, err
	}

	return NewChart(a.chartId, chartName, a.crcy, root), nil
}

func buildTreeFromDb(node tree.NodeIFace, ledgers ledgerLines, prntId uint64) (tree.NodeIFace, error) {
	var childAccounts = make(ledgerLines, 0)
	for _, line := range ledgers {
		if line.PrntId == prntId {
			childAccounts = append(childAccounts, line)
		}
	}

	for _, childAccount := range childAccounts {
		acType, ok := GetNamedAccountTypes()[childAccount.Tpe]
		if !ok {
			return nil, ErrBadAccountType
		}
		childNode := tree.NewNode(
			NewAccount(
				childAccount.Nominal,
				acType,
				childAccount.Name,
				childAccount.AcDr,
				childAccount.AcCr,
			),
			nil,
		)
		childNode, err := buildTreeFromDb(childNode, ledgers, childAccount.Id)
		if err != nil {
			return nil, err
		}
		node = node.AddChild(childNode)
	}
	return node, nil
}

//WriteTransaction writes a transaction with default datetime of Now()
func (a *Accountant) WriteTransaction(txn *SplitTransaction) (uint64, error) {
	return a.WriteTransactionWithDate(txn, time.Now())
}

//WriteTransactionWithDate writes a transaction with user supplied datetime
func (a *Accountant) WriteTransactionWithDate(txn *SplitTransaction, dt time.Time) (uint64, error) {
	if a.chartId == 0 {
		return 0, ErrNoChartId
	}

	stmt, err := a.db.Prepare("select sa_fu_add_txn(?, ?, ?, ?, ?, ?, ?, ?) as txnId")
	if err != nil {
		return 0, err
	}
	entryLen := len(txn.Entries())
	var nominals = make([]string, entryLen)
	var amounts = make([]string, entryLen)
	var tpes = make([]string, entryLen)
	acTypes := GetValuedAccountTypes()
	for i, tx := range txn.Entries() {
		nominals[i] = tx.Id().String()
		amounts[i] = fmt.Sprintf("%d", tx.Amount())
		tpes[i] = acTypes[*tx.Type()]
	}
	res, err := stmt.Query(
		a.chartId,
		txn.Note(),
		dt,
		txn.Src(),
		txn.Ref(),
		strings.Join(nominals, ","),
		strings.Join(amounts, ","),
		strings.Join(tpes, ","),
	)
	if err != nil {
		return 0, err
	}
	if res.Err() != nil {
		return 0, res.Err()
	}
	if !res.Next() {
		return 0, ErrNoJrnId
	}
	defer stmt.Close()
	var jrnId uint64
	err = res.Scan(&jrnId)
	if err != nil {
		return 0, err
	}

	return jrnId, nil
}

//FetchTransaction retrieves a journal transaction identified by its journal id
func (a *Accountant) FetchTransaction(jrnId uint64) (*SplitTransaction, error) {
	if a.chartId == 0 {
		return nil, ErrNoChartId
	}
	//the journal
	res, err := a.db.Query("select note, date, src, ref from sa_journal where id = ? and chartId = ?", jrnId, a.chartId)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	var note string
	var dt time.Time
	var src string
	var ref uint64
	defer res.Close()
	if !res.Next() {
		return nil, errors.New("cannot retrieve journal")
	}
	err = res.Scan(&note, &dt, &src, &ref)
	if err != nil {
		return nil, err
	}
	journal := NewSplitTransactionBuilder(jrnId).
		WithDate(dt).
		WithNote(note).
		WithReference(ref).
		WithSource(src)

	//journal entries
	res2, err := a.db.Query("select id, nominal, acDr, acCr from sa_journal_entry where jrnId = ?", jrnId)
	if err != nil {
		return nil, err
	}
	if res2.Err() != nil {
		return nil, res2.Err()
	}
	var id uint64
	var nominal Nominal
	var acDr, acCr uint64
	defer res2.Close()
	for res.Next() {
		err = res.Scan(&id, &nominal, &acDr, &acCr)
		if err != nil {
			return nil, err
		}
		var amount uint64
		var acType *AccountType
		if acDr == 0 {
			amount = acCr
			acType = NewAcType().Cr()
		} else {
			amount = acDr
			acType = NewAcType().Dr()
		}
		journal = journal.WithEntry(*NewEntry(nominal, amount, *acType))
	}

	return journal.Build(), nil
}

//FetchAccountJournals returns journal entries for an account
//The returned Set is a Set of SplitTransactions with only the entries for
//the required Account.  They will therefore be unbalanced.
func (a *Accountant) FetchAccountJournals(nominal Nominal) ([]*SplitTransaction, error) {
	response := make([]*SplitTransaction, 0)
	if a.chartId == 0 {
		return nil, ErrNoChartId
	}
	complexSelect := `
select j.id, j.note, j.date, j.src, j.ref, e.acDr, e.acCr
from sa_journal as j
join sa_journal_entry as e
on j.id = e.jrnid
where e.nominal = ? and j.chartId = ?
`
	res, err := a.db.Query(complexSelect, nominal.String(), a.chartId)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	var id, ref, acDr, acCr uint64
	var note, src string
	var date time.Time
	for res.Next() {
		err = res.Scan(&id, &note, &date, &src, &ref, &acDr, &acCr)
		if err != nil {
			return nil, err
		}
		var amount uint64
		var acType *AccountType
		if acDr == 0 {
			amount = acCr
			acType = NewAcType().Cr()
		} else {
			amount = acDr
			acType = NewAcType().Dr()
		}
		entry := NewEntry(nominal, amount, *acType)
		journal := NewSplitTransactionBuilder(id).
			WithNote(note).
			WithDate(date).
			WithSource(src).
			WithReference(ref).
			WithEntry(*entry).
			Build()
		response = append(response, journal)
	}

	return response, nil
}

//AddAccount adds an account (ledger) to the chart.
//Error returned if parent doesn't exist, or you try to add a second root account
func (a *Accountant) AddAccount(nominal Nominal, tpe *AccountType, name string, prnt *Nominal) error {
	var prntNominal string
	if prnt == nil {
		prntNominal = ""
	} else {
		prntNominal = prnt.String()
	}
	acTypes := GetValuedAccountTypes()
	res, err := a.db.Query("call sa_sp_add_ledger(?, ?, ?, ?, ?)",
		a.chartId,
		nominal.String(),
		acTypes[*tpe],
		name,
		prntNominal,
	)
	if err != nil {
		return err
	}
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

//DelAccount deletes an account (ledger) and all its child accounts.
//Error returned if the account has non zero debit or credit amounts
func (a *Accountant) DelAccount(nominal Nominal) error {
	res, err := a.db.Query("call sa_sp_del_ledger(?, ?)",
		a.chartId,
		nominal.String(),
	)
	if err != nil {
		return err
	}
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

//NextNominal returns the next nominal in sequence of child accounts of prnt.
//starter is given and returned if the prnt does not have child accounts
func (a *Accountant) NextNominal(prnt, starter Nominal) (*Nominal, error) {
	chart, err := a.FetchChart()
	if err != nil {
		return nil, err
	}
	res := chart.Tree().Accept(NewNodeFinder(prnt))
	if res == nil {
		return nil, ErrBadNominal
	}

	node := res.(tree.NodeIFace)
	if node.IsLeaf() {
		return &starter, nil
	}

	id, _ := strconv.ParseUint(starter.String(), 0, 64)
	for _, child := range node.GetChildren() {
		n, _ := strconv.ParseUint(child.GetValue().(*Account).Nominal().String(), 0, 64)
		if n > id {
			id = n
		}
	}
	id += 1
	next, _ := NewNominal(strconv.Itoa(int(id)))

	return &next, nil
}
