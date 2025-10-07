package utils

import (
	"encoding/json"
	"fmt"

	"github.com/bolom009/geom"
)

type Location struct {
	Polygons []*Polygon `json:"polygons"`
}

type Hole struct {
	Points   []geom.Vector2
	Viewable bool
}

type Polygon struct {
	Outer      []geom.Vector2 `json:"outer"`
	InnerHoles []Hole         `json:"innerHoles"`
	Obstacles  []Hole         `json:"obstacles"`
}

func NewLocationDataFromJSON(b []byte) (location *Location, err error) {
	err = json.Unmarshal(b, &location)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}

	return location, nil
}
