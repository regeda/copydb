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
func QueryByID(id string, resolve QueryResolve, reject QueryReject) Query {
	return &queryByID{id, resolve, reject}
}

// QueryFullScan create a query to scan all items.
func QueryFullScan(resolve QueryResolve, reject QueryReject) Query {
	return &queryFullScan{resolve, reject}
}

type queryByID struct {
	id      string
	resolve QueryResolve
	reject  QueryReject
}

func (q *queryByID) scan(ii items) {
	i, ok := ii[q.id]
	if !ok {
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
