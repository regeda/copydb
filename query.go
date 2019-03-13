package copydb

// Query scans items to perform a request.
type Query interface {
	scan(items)
}

// QueryResolve is a callback with found item.
type QueryResolve func(Item)

// QueryReject is a callback with an error.
type QueryReject func(error)

// QueryByID creates a query to get an item by identifier.
// The func accepts unix timestamp if you want to query the item not older than it.
func QueryByID(id string, unix int64, resolve QueryResolve, reject QueryReject) Query {
	return &queryByID{id, unix, resolve, reject}
}

// QueryFullScan create a query to scan all items.
func QueryFullScan(resolve QueryResolve, reject QueryReject) Query {
	return &queryFullScan{resolve, reject}
}

type queryByID struct {
	id      string
	unix    int64
	resolve QueryResolve
	reject  QueryReject
}

func (q *queryByID) scan(ii items) {
	i, ok := ii[q.id]
	if !ok || i.unix < q.unix {
		q.reject(ErrItemNotFound)
	} else {
		q.resolve(i.Item)
	}
}

type queryFullScan struct {
	resolve QueryResolve
	reject  QueryReject
}

func (q *queryFullScan) scan(ii items) {
	if ii.empty() {
		q.reject(ErrItemNotFound)
	} else {
		for _, i := range ii {
			q.resolve(i.Item)
		}
	}
}
