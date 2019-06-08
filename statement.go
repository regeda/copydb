package copydb

import (
	"context"
	"time"

	"github.com/regeda/copydb/internal/model"
)

// Statement wraps an update and execs it.
type Statement struct {
	item model.Item
}

// NewStatement creates a new statement update.
func NewStatement(id string) Statement {
	return Statement{model.Item{ID: id}}
}

// Remove erases an item.
func (s *Statement) Remove() {
	s.item.Remove = true
	s.item.Set = nil
	s.item.Unset = nil
}

// Set updates a field.
func (s *Statement) Set(name string, data []byte) {
	if s.item.Remove {
		panic("set command conflicts with remove command")
	}
	s.item.Set = append(s.item.Set, model.Item_Field{Name: name, Data: data})
}

// SetString updates a field within string.
func (s *Statement) SetString(name, data string) {
	s.Set(name, []byte(data))
}

// Unset removes a field.
func (s *Statement) Unset(name string) {
	if s.item.Remove {
		panic("unset command conflicts with remove command")
	}
	s.item.Unset = append(s.item.Unset, model.Item_Field{Name: name})
}

// Exec runs a statement.
func (s *Statement) Exec(ctx context.Context, db *DB, currtime time.Time) error {
	s.item.Unix = currtime.Unix()
	return db.exec(ctx, s.item)
}
