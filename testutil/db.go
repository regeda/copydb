package testutil

import (
	"context"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/regeda/copydb"
	"github.com/stretchr/testify/require"
)

// NewDB creates a new database instance.
func NewDB(t *testing.T, r copydb.Redis, opts ...copydb.DBOpt) *copydb.DB {
	t.Helper()

	db, err := copydb.New(r, append(opts,
		copydb.WithItemKeyPattern("test:item:{%s}"),
		copydb.WithListKey("test:items_list"),
		copydb.WithChannelKey("test:items_update"),
		copydb.WithCapacity(8),
		copydb.WithLogger(log.New(ioutil.Discard, "", 0)),
		copydb.WithPubSubOpts(copydb.PubSubOpts{
			ReceiveTimeout: 10 * time.Second,
			ChannelSize:    1,
		}),
	)...)
	require.NoError(t, err)

	return db
}

// StopDB terminates a database.
func StopDB(t *testing.T, db *copydb.DB, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	require.NoError(t, db.Stop(ctx))
}

// ServeDB runs a database.
func ServeDB(t *testing.T, db *copydb.DB) {
	t.Helper()

	go func() {
		require.NoError(t, db.Serve())
	}()

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
