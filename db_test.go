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

var now = time.Now()

func TestDB_Replicate(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db, err := copydb.New(redis)

		require.NoError(t, err)

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts := now.Add(time.Second)

		stmt := copydb.NewStatement("xxx")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitFor(waitDuration, db, "xxx", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.DefaultItem)

			assert.Equal(t, "bar", i["foo"].String())
			assert.Equal(t, "quux", i["baz"].String())
		})

		require.NoError(t, err)
	})

	t.Run("unset", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db, err := copydb.New(redis)

		require.NoError(t, err)

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts := now.Add(time.Second)

		stmt := copydb.NewStatement("yyy")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitFor(waitDuration, db, "yyy", ts.Unix())
		require.NoError(t, err)

		ts = ts.Add(time.Second)

		stmt = copydb.NewStatement("yyy")
		stmt.Unset("foo")
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitFor(waitDuration, db, "yyy", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.DefaultItem)

			_, ok := i["foo"]
			assert.False(t, ok)
			assert.Equal(t, "quux", i["baz"].String())
		})

		require.NoError(t, err)
	})

	t.Run("remove", func(t *testing.T) {
		redis := testutil.NewRedis(t)

		db, err := copydb.New(redis)

		require.NoError(t, err)

		cancel := testutil.Serve(t, db)
		defer cancel()

		ts := now.Add(time.Second)

		stmt := copydb.NewStatement("zzz")
		stmt.SetString("foo", "bar")
		stmt.SetString("baz", "quux")
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitFor(waitDuration, db, "zzz", ts.Unix())
		require.NoError(t, err)

		ts = ts.Add(time.Second)

		stmt = copydb.NewStatement("zzz")
		stmt.Remove()
		require.NoError(t, stmt.Exec(db, ts))

		err = testutil.WaitFor(waitDuration, db, "zzz", ts.Unix(), func(item copydb.Item) {
			i := item.(copydb.DefaultItem)

			assert.True(t, i.IsEmpty())
		})

		require.NoError(t, err)
	})
}
