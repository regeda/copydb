package copydb

import (
	"fmt"
)

// Pool maintains items lifecycle.
type Pool interface {
	New() Item
	Destroy(Item)
}

// Item contains functions for data manipulation.
type Item interface {
	Set(name string, data []byte)
	Unset(name string)
	Remove()
}

type defaultPool struct {
	free []DefaultItem
}

func (p *defaultPool) New() Item {
	freeNum := len(p.free)
	if freeNum > 0 {
		i := p.free[freeNum-1]
		p.free = p.free[:freeNum-1]
		return i
	}
	return make(DefaultItem)
}

func (p *defaultPool) Destroy(i Item) {
	ii, ok := i.(DefaultItem)
	if !ok {
		panic(fmt.Sprintf("default pool works with default item only, passed %T", i))
	}
	ii.Remove()
	p.free = append(p.free, ii)
}

// ItemField contains a data of a field.
type ItemField []byte

func (f ItemField) String() string {
	return string(f)
}

// DefaultItem implements Item.
type DefaultItem map[string]ItemField

// Set updates a field by name.
func (i DefaultItem) Set(name string, data []byte) {
	i[name] = ItemField(data)
}

// Unset removes a filed by name.
func (i DefaultItem) Unset(name string) {
	delete(i, name)
}

// Remove resets all fields.
func (i DefaultItem) Remove() {
	for name := range i {
		delete(i, name)
	}
}

// IsEmpty checks that the item has no fields.
func (i DefaultItem) IsEmpty() bool {
	return len(i) == 0
}
