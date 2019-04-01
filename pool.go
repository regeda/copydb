package copydb

// Pool maintains items lifecycle.
type Pool interface {
	Get() Item
	Put(Item)
}

// Item contains functions for data manipulation.
type Item interface {
	Set(name string, data []byte)
	Unset(name string)
	Remove()
}

// SimplePool implements Pool.
type SimplePool struct {
	free []Item

	New func() Item
}

// Get returns a new Item.
func (p *SimplePool) Get() Item {
	freeNum := len(p.free)
	if freeNum > 0 {
		i := p.free[freeNum-1]
		p.free = p.free[:freeNum-1]
		return i
	}
	return p.New()
}

// Put saves an item in the pool.
func (p *SimplePool) Put(i Item) {
	i.Remove()
	p.free = append(p.free, i)
}

// ItemField contains a data of a field.
type ItemField []byte

func (f ItemField) String() string {
	return string(f)
}

// SimpleItem implements Item.
type SimpleItem map[string]ItemField

// Set updates a field by name.
func (i SimpleItem) Set(name string, data []byte) {
	i[name] = ItemField(data)
}

// Unset removes a filed by name.
func (i SimpleItem) Unset(name string) {
	delete(i, name)
}

// Remove resets all fields.
func (i SimpleItem) Remove() {
	for name := range i {
		delete(i, name)
	}
}

// IsEmpty checks that the item has no fields.
func (i SimpleItem) IsEmpty() bool {
	return len(i) == 0
}
