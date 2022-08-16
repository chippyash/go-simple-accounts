package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import (
	"github.com/chippyash/go-hierarchy-tree/tree"
)

//NodeFinder finds node with required Nominal id
type NodeFinder struct {
	tree.VisitorIFace
	nom Nominal
}

//NewNodeFinder constructor
func NewNodeFinder(nominal Nominal) *NodeFinder {
	return &NodeFinder{
		nom: nominal,
	}
}

//Visit visits each mode and returns a map of account balances
func (v *NodeFinder) Visit(n tree.NodeIFace) interface{} {
	if n.GetValue().(*Account).Nominal() == v.nom {
		return n
	}
	for _, child := range n.GetChildren() {
		res := child.Accept(v)
		if res != nil {
			return res
		}
	}
	return nil
}
