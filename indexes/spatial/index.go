package spatial

import (
	"fmt"
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/google/btree"
	"github.com/regeda/copydb"
)

const earthRadiusM = 6378137.

var _ copydb.Pool = &Index{}

// Index implements Geospatial index over copydb.Pool.
type Index struct {
	level   int
	newItem func() copydb.Item
	bt      *btree.BTree
}

// NewIndex creates the index.
func NewIndex(level int, newItem func() copydb.Item) *Index {
	return &Index{
		level:   level,
		bt:      btree.New(2), // @todo figure out the proper value
		newItem: newItem,
	}
}

// New creates an item associated with the index.
func (idx *Index) New() copydb.Item {
	// @todo item's pool
	return makeItem(idx)
}

// Destroy the item from the index.
func (idx *Index) Destroy(it copydb.Item) {
	i, ok := it.(*item)
	if !ok {
		panic(fmt.Sprintf("spatial index works with spatial item only, passed %T", it))
	}
	if i.idx != idx {
		panic("item was created by another index")
	}

	i.Remove()
}

func (idx *Index) remove(it *item) {
	if !it.isIndexed() {
		return
	}

	bitem := idx.bt.Get(makeNode(it.cellID))
	if bitem != nil {
		n := bitem.(*node)
		n.remove(it)
		if n.isEmpty() {
			idx.bt.Delete(bitem)
		}
	}
}

func (idx *Index) move(it *item, lon, lat float64) (s2.LatLng, s2.CellID) {
	idx.remove(it)
	return idx.add(it, lon, lat)
}

func (idx *Index) add(it *item, lon, lat float64) (s2.LatLng, s2.CellID) {
	latlng := s2.LatLngFromDegrees(lat, lon)
	cellID := s2.CellIDFromLatLng(latlng)
	cellIDOnStorageLevel := cellID.Parent(idx.level)

	n := makeNode(cellIDOnStorageLevel)

	bitem := idx.bt.Get(n)
	if bitem != nil {
		n = bitem.(*node)
	}

	n.add(it)

	idx.bt.ReplaceOrInsert(n)

	return latlng, cellIDOnStorageLevel
}

// @todo introduce max results and filter function
func (idx *Index) search(req SearchRequest, resolve copydb.QueryResolve, reject copydb.QueryReject) {
	latlng := s2.LatLngFromDegrees(req.Lat, req.Lon)
	centerAngle := s1.Angle(req.Radius / earthRadiusM)

	cap := s2.CapFromCenterAngle(s2.PointFromLatLng(latlng), centerAngle)
	rc := s2.RegionCoverer{MaxLevel: idx.level}
	cu := rc.Covering(cap)

	visitor := visitor{
		angle:   centerAngle,
		latlng:  latlng,
		closest: newItemsQueue(req.Limit),
		k:       req.Limit,
	}

	if visitor.k == 0 {
		visitor.k = int(math.MaxInt32)
	}

	for _, cellID := range cu {
		if cellID.Level() < idx.level {
			begin := cellID.ChildBeginAtLevel(idx.level)
			end := cellID.ChildEndAtLevel(idx.level)
			idx.bt.AscendRange(makeNode(begin), makeNode(end.Next()), func(bitem btree.Item) bool {
				visitor.visit(bitem.(*node).items)
				return true
			})
		} else {
			bitem := idx.bt.Get(makeNode(cellID))
			if bitem != nil {
				visitor.visit(bitem.(*node).items)
			}
		}
	}

	if visitor.closest.Len() == 0 {
		reject(copydb.ErrItemNotFound)
	} else {
		for _, v := range visitor.closest {
			resolve(v.Item)
		}
	}
}

// SearchRequest contains parameters for nearby search.
type SearchRequest struct {
	Lon, Lat float64
	Radius   float64
	Limit    int
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
