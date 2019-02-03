package providers

import (
	"time"

	geo "github.com/paulmach/go.geo"
)

type Coordinate struct {
	Point geo.Point
	Ts    time.Time
}

type Provider struct {
	ID     string
	Coord  Coordinate
	Status string
}

func NewProvider(id string) *Provider {
	return &Provider{ID: id}
}
