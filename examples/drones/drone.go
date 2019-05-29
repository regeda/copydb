package drones

import (
	geojson "github.com/paulmach/go.geojson"
)

type Drone struct {
	ID     string            `json:"id"`
	Geom   *geojson.Geometry `json:"geom"`
	Status string            `json:"status"`
}

func NewDrone(id string) *Drone {
	return &Drone{ID: id}
}
