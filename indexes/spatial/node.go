package spatial

import (
	"github.com/golang/geo/s2"
	"github.com/google/btree"
)

var _ btree.Item = &node{}

type node struct {
	cellID s2.CellID
	items  []*item
}

func makeNode(cellID s2.CellID) *node {
	return &node{cellID: cellID}
}

func (n *node) Less(than btree.Item) bool {
	return n.cellID < than.(*node).cellID
}

func (n *node) add(it *item) {
	n.items = append(n.items, it)
}

func (n *node) remove(it *item) {
	i, ok := n.indexOf(it)
	if !ok {
		return
	}
	copy(n.items[i:], n.items[i+1:])
	n.items[len(n.items)-1] = nil
	n.items = n.items[:len(n.items)-1]
}

func (n *node) isEmpty() bool {
	return len(n.items) == 0
}

func (n *node) indexOf(it *item) (int, bool) {
	for i, ii := range n.items {
		if ii == it {
			return i, true
		}
	}
	return 0, false
}
