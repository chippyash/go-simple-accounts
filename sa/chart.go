package sa

/**
 * Simple Double Entry Accounting V3 for Go

 * @author Ashley Kitson
 * @copyright Ashley Kitson, 2022, UK
 * @license BSD-3-Clause See LICENSE.md
 */

import "github.com/chippyash/go-hierarchy-tree/tree"

//Chart is a COA
type Chart struct {
	id   uint64
	tree tree.NodeIFace
	name string
	crcy string
}

//NewChart Constructor. Returns a new COA
func NewChart(id uint64, name string, crcy string, chartTree tree.NodeIFace) *Chart {
	if chartTree == nil {
		chartTree = tree.NewNode(nil, nil)
	}
	return &Chart{
		id:   id,
		tree: chartTree,
		name: name,
		crcy: crcy,
	}
}

func (c *Chart) findAccountNode(nominal Nominal) tree.NodeIFace {
	filterFunc := func(n tree.NodeIFace) bool {
		return n.GetValue().(*Account).Nominal() == nominal
	}
	filter := tree.NewFilterVisitor(filterFunc)
	res := c.tree.Accept(filter).([]tree.NodeIFace)
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

//GetAccount returns first account matching nominal. Can return nil if not found
func (c *Chart) GetAccount(nominal Nominal) *Account {
	res := c.findAccountNode(nominal)
	if res == nil {
		return nil
	}
	return res.GetValue().(*Account)
}

func (c *Chart) findAccountName(name string) tree.NodeIFace {
	filterFunc := func(n tree.NodeIFace) bool {
		return n.GetValue().(*Account).Name() == name
	}
	filter := tree.NewFilterVisitor(filterFunc)
	res := c.tree.Accept(filter).([]tree.NodeIFace)
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

//GetAccountByName returns first account matching nominal. Can return nil if not found
func (c *Chart) GetAccountByName(name string) *Account {
	res := c.findAccountName(name)
	if res == nil {
		return nil
	}
	return res.GetValue().(*Account)
}

//HasAccount tests if chart contains an account matching nominal code
func (c *Chart) HasAccount(nominal Nominal) bool {
	return c.GetAccount(nominal) != nil
}

//GetParentId returns the parent nominal code for the required account. Returned code can be empty
//indicating that the account has no parent, and thus is likely to be the root account
func (c *Chart) GetParentId(nominal Nominal) Nominal {
	res := c.findAccountNode(nominal)
	if res == nil {
		return ""
	}
	return res.GetParent().GetValue().(*Account).Nominal()
}

//Tree returns the COA root tree node for direct manipulation
func (c *Chart) Tree() tree.NodeIFace {
	return c.tree
}

//Name returns the COA name
func (c *Chart) Name() string {
	return c.name
}

//Id returns the chart storage id
func (c *Chart) Id() uint64 {
	return c.id
}

//Crcy returns the COA currency (3 char crcy code)
func (c *Chart) Crcy() string {
	return c.crcy
}

//SetRootNode directly sets the COA root node
func (c *Chart) SetRootNode(root tree.NodeIFace) *Chart {
	c.tree = root
	return c
}
