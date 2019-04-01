package spatial_test

import (
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
			Radius: 1e5,
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

	// nearest
	a.Set("foo", []byte("bar"))
	a.Set("geom", []byte(`{"type":"Point","coordinates":[-73.98841112852097,40.73780734529185]}`))

	// farest
	b.Set("baz", []byte("quux"))
	b.Set("geom", []byte(`{"type":"Point","coordinates":[-74.01154518127441,40.71486662585682]}`))

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

			assert.Equal(t, "bar", it["foo"].String())
		},
		func(err error) {
			assert.FailNow(t, "no error should be", "error was triggered: %v", err)
		},
	)

	query(idx)

	require.Equal(t, 1, found)
}

func TestIndex_NothingFoundAfterItemDestroyed(t *testing.T) {
	idx := spatial.NewIndex(13, newItem)

	a := idx.Get()

	a.Set("foo", []byte("bar"))
	a.Set("geom", []byte(`{"type":"Point","coordinates":[-73.98841112852097,40.73780734529185]}`))

	idx.Put(a)

	var err error

	query := spatial.Search(
		spatial.SearchRequest{
			Lon:    -73.98914,
			Lat:    40.73769,
			Radius: 1e5,
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

	a.Set("foo", []byte("bar"))
	a.Set("geom", []byte(`{"type":"Point","coordinates":[-73.98841112852097,40.73780734529185]}`))

	a.Unset("geom")

	var err error

	query := spatial.Search(
		spatial.SearchRequest{
			Lon:    -73.98914,
			Lat:    40.73769,
			Radius: 1e5,
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
