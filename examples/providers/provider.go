package providers

import (
	geojson "github.com/paulmach/go.geojson"
)

type Provider struct {
	ID     string            `json:"id"`
	Geom   *geojson.Geometry `json:"geom"`
	Status string            `json:"status"`
}

func NewProvider(id string) *Provider {
	return &Provider{ID: id}
}
