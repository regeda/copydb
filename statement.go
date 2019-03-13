package copydb

import "time"

// Statement wraps an update and execs it.
type Statement struct {
	upd *Update
}

// NewStatement creates a new statement update.
func NewStatement(id string) Statement {
	return Statement{allocUpdate(id)}
}

// Remove erases an item.
func (a *Statement) Remove() {
	a.upd.Remove = true
	a.upd.Set = nil
	a.upd.Unset = nil
}

// Set updates a field.
func (a *Statement) Set(name string, data []byte) {
	if a.upd.Remove {
		panic("set command conflicts with remove command")
	}
	a.upd.Set = append(a.upd.Set, Update_Field{Name: name, Data: data})
}

// SetString updates a field within string.
func (a *Statement) SetString(name, data string) {
	a.Set(name, []byte(data))
}

// Unset removes a field.
func (a *Statement) Unset(name string) {
	if a.upd.Remove {
		panic("unset command conflicts with remove command")
	}
	a.upd.Unset = append(a.upd.Unset, Update_Field{Name: name})
}

// Exec runs a statement.
func (a *Statement) Exec(db *DB, currtime time.Time) error {
	defer freeUpdate(a.upd)
	a.upd.Unix = currtime.Unix()
	return db.replicate(a.upd)
}
