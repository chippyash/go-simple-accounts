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
	"github.com/chippyash/go-hierarchy-tree/tree"
)

//NodeSaver saves account ledger definitions to the database
type NodeSaver struct {
	tree.VisitorIFace
	db *sql.DB
	id uint64
}

//NewNodeSaver constructor
func NewNodeSaver(db *sql.DB) *NodeSaver {
	return &NodeSaver{
		db: db,
	}
}

//Visit store each tree node in the DB, Returns error or nil
func (v *NodeSaver) Visit(n tree.NodeIFace) interface{} {
	currAc := n.GetValue().(*Account)
	tpe, ok := GetValuedAccountTypes()[*currAc.Type()]
	if !ok {
		return errors.New("cannot retrieve account type string from value")
	}
	var prntNominal string
	if n.IsRoot() {
		prntNominal = ""
	} else {
		prntNominal = n.GetParent().GetValue().(*Account).Nominal().String()
	}

	_, err := v.db.Exec("call sa_sp_add_ledger(?, ?, ?, ?, ?)", currAc.chartId, currAc.Nominal().String(), tpe, currAc.Name(), prntNominal)
	if err != nil {
		return err
	}

	for _, child := range n.GetChildren() {
		errv := child.Accept(v)
		if errv != nil {
			return errv
		}
	}

	return nil
}
