package spatial_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/regeda/copydb"
	"github.com/regeda/copydb/indexes/spatial"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newItem() copydb.Item {
	return make(copydb.SimpleItem)
}

func TestIndex_NothingFoundIfIndexEmpty(t *testing.T) {
	idx := spatial.NewIndex(13, newItem)

	var err error

	query := spatial.Search(
		spatial.SearchRequest{
			Lon:    -73.98914,
			Lat:    40.73769,
			Radius: 10000,
		},
		func(copydb.Item) {
			assert.FailNow(t, "no items should be found")
		},
		func(e error) {
			err = e
		},
	)

	query(idx)

	assert.EqualError(t, err, "item not found")
}

func TestIndex_FoundNearest(t *testing.T) {
	idx := spatial.NewIndex(13, newItem)

	a := idx.Get()
	b := idx.Get()

	a.Set("name", []byte("nearest"))
	a.Set("geom", []byte(`{"type":"Point","coordinates":[-73.98841112852097,40.73780734529185]}`))

	b.Set("name", []byte("farest"))
	b.Set("geom", []byte(`{"type":"Point","coordinates":[-74.01154518127441,40.71486662585682]}`))

	t.Run("in radius", func(t *testing.T) {
		var found int

		query := spatial.Search(
			spatial.SearchRequest{
				Lon:    -73.98914,
				Lat:    40.73769,
				Radius: 100,
			},
			func(item copydb.Item) {
				found++

				it := item.(copydb.SimpleItem)

				assert.Equal(t, "nearest", it["name"].String())
			},
			func(err error) {
				assert.FailNow(t, "no error should be", "error was triggered: %v", err)
			},
		)

		query(idx)

		require.Equal(t, 1, found)
	})

	t.Run("with filter", func(t *testing.T) {
		var found int

		query := spatial.Search(
			spatial.SearchRequest{
				Lon:    -73.98914,
				Lat:    40.73769,
				Radius: 10000,
				Filter: func(item copydb.Item) bool {
					it := item.(copydb.SimpleItem)
					return it["name"].String() == "farest"
				},
			},
			func(item copydb.Item) {
				found++

				it := item.(copydb.SimpleItem)

				assert.Equal(t, "farest", it["name"].String())
			},
			func(err error) {
				assert.FailNow(t, "no error should be", "error was triggered: %v", err)
			},
		)

		query(idx)

		require.Equal(t, 1, found)
	})

}

func TestIndex_NothingFoundAfterItemReturnedInPool(t *testing.T) {
	idx := spatial.NewIndex(13, newItem)

	a := idx.Get()

	a.Set("geom", []byte(`{"type":"Point","coordinates":[-73.98841112852097,40.73780734529185]}`))

	idx.Put(a)

	var err error

	query := spatial.Search(
		spatial.SearchRequest{
			Lon:    -73.98914,
			Lat:    40.73769,
			Radius: 10000,
		},
		func(copydb.Item) {
			assert.FailNow(t, "no items should be found")
		},
		func(e error) {
			err = e
		},
	)

	query(idx)

	assert.EqualError(t, err, "item not found")
}

func TestIndex_NothingFoundIfGeomRemoved(t *testing.T) {
	idx := spatial.NewIndex(13, newItem)

	a := idx.Get()

	a.Set("geom", []byte(`{"type":"Point","coordinates":[-73.98841112852097,40.73780734529185]}`))

	a.Unset("geom")

	var err error

	query := spatial.Search(
		spatial.SearchRequest{
			Lon:    -73.98914,
			Lat:    40.73769,
			Radius: 2000,
		},
		func(copydb.Item) {
			assert.FailNow(t, "no items should be found")
		},
		func(e error) {
			err = e
		},
	)

	query(idx)

	assert.EqualError(t, err, "item not found")
}

// BenchmarkIndex_Search-4            20000            139719 ns/op           11044 B/op         41 allocs/op
func BenchmarkIndex_Search(b *testing.B) {
	idx := spatial.NewIndex(13, newItem)

	for i := 0; i < 2000; i++ {
		a := idx.Get()
		a.Set("geom", []byte(fmt.Sprintf(`{"type":"Point","coordinates":[%g,%g]}`, -73+rand.Float64(), 40+rand.Float64())))
	}

	query := spatial.Search(
		spatial.SearchRequest{
			Lon:    -73 + rand.Float64(),
			Lat:    40 + rand.Float64(),
			Radius: 10000,
			Limit:  10,
		},
		func(copydb.Item) {
		},
		func(error) {
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query(idx)
	}
}
