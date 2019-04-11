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
		defer db.MustStop(waitDuration)

		testutil.Serve(db)

		ts := time.Now()

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err := testutil.WaitForItem(waitDuration, db, "xxx", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.SimpleItem)

			assert.Equal(t, "bar", i["foo"].String())
			assert.Equal(t, "quux", i["baz"].String())
		})

		require.NoError(t, err)
	})

	t.Run("unset", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db := copydb.MustNew(redis)
		defer db.MustStop(waitDuration)

		testutil.Serve(db)

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
			i := item.(copydb.SimpleItem)

			_, ok := i["foo"]
			assert.False(t, ok)
			assert.Equal(t, "quux", i["baz"].String())
		})

		require.NoError(t, err)
	})

	t.Run("remove", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db := copydb.MustNew(redis)
		defer db.MustStop(waitDuration)

		testutil.Serve(db)

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
			i := item.(copydb.SimpleItem)

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
		defer db.MustStop(waitDuration)

		testutil.Serve(db)

		ts1 := time.Now()

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		require.NoError(t, stmt.Exec(db, ts1))

		ts2 := ts1.Add(ttl * 2)

		stmt = copydb.NewStatement("yyy")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts2))

		require.NoError(t,
			testutil.WaitForItem(ttl, db, "xxx", ts1.Unix()),
		)
		require.NoError(t,
			testutil.WaitForItem(ttl, db, "yyy", ts2.Unix()),
		)

		require.EqualError(t,
			testutil.WaitForError(ttl+time.Second, db, "xxx", ts1.Unix()),
			"item not found",
		)
		require.NoError(t,
			testutil.WaitForItem(ttl, db, "yyy", ts2.Unix()),
		)
	})

	t.Run("on_start", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db1 := copydb.MustNew(redis, copydb.WithTTL(ttl))
		defer db1.MustStop(waitDuration)

		testutil.Serve(db1)

		ts1 := time.Now()

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		require.NoError(t, stmt.Exec(db1, ts1))

		ts2 := ts1.Add(ttl * 2)

		stmt = copydb.NewStatement("yyy")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db1, ts2))

		require.NoError(t,
			testutil.WaitForItem(ttl, db1, "xxx", ts1.Unix()),
		)
		require.NoError(t,
			testutil.WaitForItem(ttl, db1, "yyy", ts2.Unix()),
		)

		time.Sleep(ttl + time.Second)

		db2 := copydb.MustNew(redis, copydb.WithTTL(ttl))
		defer db2.MustStop(waitDuration)

		testutil.Serve(db2)

		require.EqualError(t,
			testutil.WaitForError(ttl, db2, "xxx", ts1.Unix()),
			"item not found",
		)
		require.NoError(t,
			testutil.WaitForItem(ttl, db2, "yyy", ts2.Unix()),
		)

	})
}

func TestDB_ResolveVersionConflict(t *testing.T) {
	redis := testutil.NewRedis(t)

	db1 := copydb.MustNew(redis)
	defer db1.MustStop(waitDuration)

	testutil.Serve(db1)

	db2 := copydb.MustNew(redis)
	defer db2.MustStop(waitDuration)

	testutil.Serve(db2)

	ts1 := time.Now()

	stmt1 := copydb.NewStatement("xxx")
	stmt1.SetString("foo", "bar")
	require.NoError(t, stmt1.Exec(db1, ts1))

	require.NoError(t,
		testutil.WaitForItem(waitDuration, db1, "xxx", ts1.Unix()),
	)

	require.NoError(t,
		testutil.WaitForItem(waitDuration, db2, "xxx", ts1.Unix()),
	)

	require.NoError(t, db2.Stop(waitDuration))

	ts2 := ts1.Add(time.Second)

	stmt2 := copydb.NewStatement("xxx")
	stmt2.SetString("baz", "quux")
	require.NoError(t, stmt2.Exec(db1, ts2))

	require.NoError(t,
		testutil.WaitForItem(waitDuration, db1, "xxx", ts2.Unix()),
	)

	testutil.Serve(db2)

	require.EqualError(t,
		testutil.WaitForError(waitDuration, db2, "xxx", ts2.Unix()),
		"item not found",
	)

	ts3 := ts2.Add(time.Second)

	stmt3 := copydb.NewStatement("xxx")
	stmt3.Unset("foo")
	require.NoError(t, stmt3.Exec(db1, ts3))

	require.NoError(t,
		testutil.WaitForItem(waitDuration, db1, "xxx", ts3.Unix()),
	)

	require.NoError(t,
		testutil.WaitForItem(waitDuration, db2, "xxx", ts3.Unix()),
	)

}
