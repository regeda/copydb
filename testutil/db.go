package testutil

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/regeda/copydb"
)

// Serve runs DB.Serve with cancellation func.
func Serve(t *testing.T, db *copydb.DB) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := db.Serve(ctx); err != nil {
			t.Fatalf("serve failed: %v", err)
		}
	}()

	// wait until queries become acceptable
	db.Queries() <- copydb.QueryFullScan(func(copydb.Item) {}, func(error) {})

	return cancel
}

// WaitFor waits until an identifier will be found in the database.
func WaitFor(d time.Duration, db *copydb.DB, id string, unix int64, resolve ...copydb.QueryResolve) error {
	doneCh := make(chan struct{})
	defer close(doneCh)

	foundCh := make(chan struct{})

	go func() {
		errCh := make(chan error)
		for {
			select {
			case db.Queries() <- copydb.QueryByID(id, unix,
				func(item copydb.Item) {
					defer close(foundCh)
					for _, r := range resolve {
						r(item)
					}
				},
				func(err error) {
					errCh <- err
				},
			):
			case <-doneCh:
				return
			}

			select {
			case <-errCh:
			case <-doneCh:
				return
			}
		}
	}()

	select {
	case <-foundCh:
		return nil
	case <-time.Tick(d):
		return errors.New("timeout exceeded")
	}
}
