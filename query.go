package copydb

// Query scans items to perform a request.
type Query interface {
	scan(*DB)
}

type queryFunc func(*DB)

func (f queryFunc) scan(db *DB) {
	f(db)
}

// QueryResolve is a callback with found item.
type QueryResolve func(Item)

// QueryReject is a callback with an error.
type QueryReject func(error)

// QueryByID creates a query to get an item by identifier.
// The func accepts unix timestamp if you want to query the item not older than it.
func QueryByID(id string, unix int64, resolve QueryResolve, reject QueryReject) Query {
	return queryFunc(func(db *DB) {
		i, ok := db.items[id]
		if ok && i.unix >= unix {
			resolve(i.Item)
		} else {
			reject(ErrItemNotFound)
		}
	})
}

// QueryAll creates a query to scan all items.
func QueryAll(resolve QueryResolve, reject QueryReject) Query {
	return queryFunc(func(db *DB) {
		if db.items.empty() {
			reject(ErrItemNotFound)
		} else {
			for _, i := range db.items {
				resolve(i.Item)
			}
		}
	})
}

// QueryPool creates a query to scan the pool directly.
func QueryPool(scan func(Pool)) Query {
	return queryFunc(func(db *DB) {
		scan(db.pool)
	})
}

// QueryStats creates a query to scan db statistics.
func QueryStats(scan func(Stats)) Query {
	return queryFunc(func(db *DB) {
		scan(db.stats)
	})
}
