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
	free []defaultItem
}

func (p *defaultPool) New() Item {
	freeNum := len(p.free)
	if freeNum > 0 {
		i := p.free[freeNum-1]
		p.free = p.free[:freeNum-1]
		return i
	}
	return make(defaultItem)
}

func (p *defaultPool) Destroy(i Item) {
	ii, ok := i.(defaultItem)
	if !ok {
		panic(fmt.Sprintf("default pool works with default item only, passed %T", i))
	}
	ii.Remove()
	p.free = append(p.free, ii)
}

type defaultItem map[string][]byte

func (i defaultItem) Set(name string, data []byte) {
	i[name] = data
}

func (i defaultItem) Unset(name string) {
	delete(i, name)
}

func (i defaultItem) Remove() {
	for name := range i {
		delete(i, name)
	}
}
