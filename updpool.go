package copydb

import "sync"

var updpool = sync.Pool{
	New: func() interface{} {
		return new(Update)
	},
}

func allocUpdate(id string) *Update {
	u := updpool.Get().(*Update)
	u.ID = id
	return u
}

func freeUpdate(u *Update) {
	u.Reset()
	updpool.Put(u)
}
