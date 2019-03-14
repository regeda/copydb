package copydb_test

import (
	"testing"
	"time"

	"github.com/regeda/copydb"
	"github.com/regeda/copydb/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const waitDuration = 5 * time.Second

func TestDB_Replicate(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db := copydb.MustNew(redis)

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts := time.Now()

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err := testutil.WaitForItem(waitDuration, db, "xxx", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.DefaultItem)

			assert.Equal(t, "bar", i["foo"].String())
			assert.Equal(t, "quux", i["baz"].String())
		})

		require.NoError(t, err)
	})

	t.Run("unset", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db := copydb.MustNew(redis)

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts := time.Now()

		stmt := copydb.NewStatement("yyy")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err := testutil.WaitForItem(waitDuration, db, "yyy", ts.Unix())
		require.NoError(t, err)

		ts = ts.Add(time.Second)

		stmt = copydb.NewStatement("yyy")
		stmt.Unset("foo")
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitForItem(waitDuration, db, "yyy", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.DefaultItem)

			_, ok := i["foo"]
			assert.False(t, ok)
			assert.Equal(t, "quux", i["baz"].String())
		})

		require.NoError(t, err)
	})

	t.Run("remove", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db := copydb.MustNew(redis)

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts := time.Now()

		stmt := copydb.NewStatement("zzz")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err := testutil.WaitForItem(waitDuration, db, "zzz", ts.Unix())
		require.NoError(t, err)

		ts = ts.Add(time.Second)

		stmt = copydb.NewStatement("zzz")
		stmt.Remove()
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitForItem(waitDuration, db, "zzz", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.DefaultItem)

			assert.True(t, i.IsEmpty())
		})

		require.NoError(t, err)
	})
}

func TestDB_EvictExpired(t *testing.T) {
	ttl := 5 * time.Second

	t.Run("by_timer", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db := copydb.MustNew(redis, copydb.WithTTL(ttl))

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts1 := time.Now()

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		require.NoError(t, stmt.Exec(db, ts1))

		ts2 := ts1.Add(ttl * 2)

		stmt = copydb.NewStatement("yyy")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts2))

		require.NoError(t,
			testutil.WaitForItem(waitDuration, db, "xxx", ts1.Unix()),
		)
		require.NoError(t,
			testutil.WaitForItem(waitDuration, db, "yyy", ts2.Unix()),
		)

		require.EqualError(t,
			testutil.WaitForError(waitDuration, db, "xxx", ts1.Unix()),
			"item not found",
		)

		require.NoError(t,
			testutil.WaitForItem(waitDuration, db, "yyy", ts2.Unix()),
		)
	})

	t.Run("on_start", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db1 := copydb.MustNew(redis, copydb.WithTTL(ttl))

		cancel1 := testutil.Serve(t, db1)
		defer cancel1()

		ts := time.Now()

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		require.NoError(t, stmt.Exec(db1, ts))

		require.NoError(t,
			testutil.WaitForItem(waitDuration, db1, "xxx", ts.Unix()),
		)

		time.Sleep(ttl + time.Second)

		db2 := copydb.MustNew(redis, copydb.WithTTL(ttl))

		cancel2 := testutil.Serve(t, db2)
		defer cancel2()

		require.EqualError(t,
			testutil.WaitForError(waitDuration, db2, "xxx", ts.Unix()),
			"item not found",
		)
	})
}
