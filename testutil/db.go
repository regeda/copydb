package testutil

import (
	"context"
	"time"

	"github.com/regeda/copydb"
)

// Serve runs DB.Serve with cancellation func.
func Serve(db *copydb.DB) {
	go db.MustServe()

	// wait until queries become acceptable
	db.QueriesIn() <- copydb.QueryAll(func(copydb.Item) {}, func(error) {})
}

// WaitForItem waits until an identifier will be found in the database.
func WaitForItem(d time.Duration, db *copydb.DB, id string, unix int64, resolve ...copydb.QueryResolve) error {
	ctx, cancel := context.WithTimeout(context.TODO(), d)
	defer cancel()

	var found bool

	for !found {
		if err := db.Query(ctx, copydb.QueryByID(id, unix,
			func(item copydb.Item) {
				for _, r := range resolve {
					r(item)
				}
				found = true
			},
			func(error) {
			},
		)); err != nil {
			return err
		}
	}

	return nil
}

// WaitForError waits until an error will be returned.
func WaitForError(d time.Duration, db *copydb.DB, id string, unix int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), d)
	defer cancel()

	var err error

	for err == nil {
		if qerr := db.Query(ctx, copydb.QueryByID(id, unix,
			func(copydb.Item) {
			},
			func(e error) {
				err = e
			},
		)); qerr != nil {
			return nil
		}
	}

	return err
}
