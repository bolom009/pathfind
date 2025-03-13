package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bolom009/geom"
)

type Canvas struct {
	W float32 `json:"w"`
	H float32 `json:"h"`
}

type polygonToolJSON struct {
	Canvas   Canvas           `json:"canvas"`
	Polygons [][]geom.Vector2 `json:"polygons"`
}

// NewPolygonsFromJSON loads a polygon set from JSON data as exported by the
// [Polygon Constructor]: https://alaricus.github.io/PolygonConstructor/
func NewPolygonsFromJSON(jsonData []byte) (polygon []geom.Vector2, holes [][]geom.Vector2, size geom.Vector2, err error) {
	var jsonStruct polygonToolJSON
	err = json.Unmarshal(jsonData, &jsonStruct)
	if err != nil {
		return nil, nil, geom.Vector2{}, fmt.Errorf("could not unmarshal polygon JSON: %w", err)
	}

	polygon = make([]geom.Vector2, 0)
	for _, pp := range jsonStruct.Polygons[0] {
		polygon = append(polygon, geom.Vector2{X: pp.X, Y: pp.Y})
	}

	// Example hole: a small square hole in the middle: (4,4), (6,4), (6,6), (4,6)
	holes = make([][]geom.Vector2, 0)
	for _, pp := range jsonStruct.Polygons[1:] {
		hole := make([]geom.Vector2, len(pp))
		for j, pp := range pp {
			hole[j] = geom.Vector2{X: pp.X, Y: pp.Y}
		}

		holes = append(holes, hole)
	}

	return polygon, holes, jsonStruct.size(), nil
}

func (pj *polygonToolJSON) size() geom.Vector2 {
	return geom.Vector2{X: pj.Canvas.W, Y: pj.Canvas.H}
}
