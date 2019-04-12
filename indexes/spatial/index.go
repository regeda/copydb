package spatial

import (
	"fmt"
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/regeda/copydb"
)

const earthRadiusM = 6371010.

var _ copydb.Pool = &Index{}

// Index implements Geospatial index over copydb.Pool.
type Index struct {
	*copydb.SimplePool

	level   int
	grid    map[s2.CellID]cell
	visitor visitor
}

// NewIndex creates the index.
func NewIndex(level int, newItem func() copydb.Item) *Index {
	idx := &Index{
		level: level,
		grid:  make(map[s2.CellID]cell),
	}

	idx.SimplePool = &copydb.SimplePool{
		New: func() copydb.Item {
			return makeItem(idx, newItem())
		},
	}

	return idx
}

func (idx *Index) remove(it *item) {
	if c, ok := idx.grid[it.cellID]; ok {
		c.remove(it)
	}
}

func (idx *Index) move(it *item, lon, lat float64) {
	latlng := s2.LatLngFromDegrees(lat, lon)
	cellID := s2.CellIDFromLatLng(latlng).Parent(idx.level)

	if it.cellID != cellID {
		idx.remove(it)

		c, ok := idx.grid[cellID]
		if !ok {
			c = make(cell)
			idx.grid[cellID] = c
		}

		c.add(it)
	}

	it.latlng = latlng
	it.cellID = cellID
}

func (idx *Index) search(req SearchRequest, resolve copydb.QueryResolve, reject copydb.QueryReject) {
	idx.visitor.angle = s1.Angle(req.Radius / earthRadiusM)
	idx.visitor.latlng = s2.LatLngFromDegrees(req.Lat, req.Lon)
	idx.visitor.k = req.Limit
	idx.visitor.filter = req.Filter

	if idx.visitor.k == 0 {
		idx.visitor.k = int(math.MaxInt32)
	}

	cellID := s2.CellIDFromLatLng(idx.visitor.latlng).Parent(idx.level)

	visited := make(map[s2.CellID]bool)
	walk := []s2.CellID{cellID}

	for cursor := 0; cursor < len(walk); cursor++ {
		cellID = walk[cursor]

		if visited[cellID] {
			continue
		}

		if c, ok := idx.grid[cellID]; ok {
			idx.visitor.visit(c)
		}

		visited[cellID] = true

		for _, cellID := range cellID.EdgeNeighbors() {
			if !visited[cellID] && cellID.LatLng().Distance(idx.visitor.latlng) < idx.visitor.angle {
				walk = append(walk, cellID)
			}
		}
	}

	if idx.visitor.closest.Len() == 0 {
		reject(copydb.ErrItemNotFound)
	} else {
		for _, v := range idx.visitor.closest {
			resolve(v.Item)
		}
		idx.visitor.closest.reset()
	}
}

// SearchRequest contains parameters for nearby search.
type SearchRequest struct {
	Lon, Lat float64
	Radius   float64
	Limit    int
	Filter   func(copydb.Item) bool
}

// Search returns items by a coordinate and a radius.
func Search(req SearchRequest, resolve copydb.QueryResolve, reject copydb.QueryReject) func(copydb.Pool) {
	return func(pool copydb.Pool) {
		idx, ok := pool.(*Index)
		if !ok {
			panic(fmt.Sprintf("spatial search works with spatial index only, given %T", pool))
		}
		idx.search(req, resolve, reject)
	}
}
