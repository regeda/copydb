package testutil

import (
	"errors"
	"time"

	"github.com/regeda/copydb"
)

// Serve runs DB.Serve with cancellation func.
func Serve(db *copydb.DB) {
	go db.MustServe()

	// wait until queries become acceptable
	db.Queries() <- copydb.QueryFullScan(func(copydb.Item) {}, func(error) {})
}

// WaitForItem waits until an identifier will be found in the database.
func WaitForItem(d time.Duration, db *copydb.DB, id string, unix int64, resolve ...copydb.QueryResolve) error {
	doneCh := make(chan struct{})
	defer close(doneCh)

	foundCh := make(chan struct{})

	go func() {
		errCh := make(chan struct{})
		for {
			select {
			case db.Queries() <- copydb.QueryByID(id, unix,
				func(item copydb.Item) {
					defer close(foundCh)
					for _, r := range resolve {
						r(item)
					}
				},
				func(error) {
					errCh <- struct{}{}
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

// WaitForError waits until an error will be returned.
func WaitForError(d time.Duration, db *copydb.DB, id string, unix int64) error {
	doneCh := make(chan struct{})
	defer close(doneCh)

	errCh := make(chan error)

	go func() {
		resolveCh := make(chan struct{})
		for {
			select {
			case db.Queries() <- copydb.QueryByID(id, unix,
				func(copydb.Item) {
					resolveCh <- struct{}{}
				},
				func(err error) {
					errCh <- err
				},
			):
			case <-doneCh:
				return
			}

			select {
			case <-resolveCh:
			case <-doneCh:
				return
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.Tick(d):
		return nil
	}
}
