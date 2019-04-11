package spatial

import (
	"github.com/golang/geo/s2"
	geojson "github.com/paulmach/go.geojson"
	"github.com/regeda/copydb"
)

var _ copydb.Item = &item{}

type item struct {
	it     copydb.Item
	idx    *Index
	cellID s2.CellID
	latlng s2.LatLng
}

func makeItem(idx *Index, it copydb.Item) *item {
	return &item{
		idx:    idx,
		it:     it,
		cellID: s2.SentinelCellID,
	}
}

func (i *item) Set(name string, data []byte) {
	switch name {
	case "geom":
		geom, err := geojson.UnmarshalGeometry(data)
		if err == nil && geom.IsPoint() {
			i.idx.move(i, geom.Point[0], geom.Point[1])
		}
	}
	i.it.Set(name, data)
}

func (i *item) detachFromS2() {
	i.idx.remove(i)
	i.cellID = s2.SentinelCellID
	i.latlng = s2.LatLng{}
}

func (i *item) Unset(name string) {
	switch name {
	case "geom":
		i.detachFromS2()
	}
	i.it.Unset(name)
}

func (i *item) Remove() {
	i.detachFromS2()
	i.it.Remove()
}
