package copydb

import (
	"sync"

	"github.com/regeda/copydb/internal/model"
)

var updpool = sync.Pool{
	New: func() interface{} {
		return new(model.Update)
	},
}

func allocUpdate(id string) *model.Update {
	u := updpool.Get().(*model.Update)
	u.ID = id
	return u
}

func freeUpdate(u *model.Update) {
	u.Reset()
	updpool.Put(u)
}
