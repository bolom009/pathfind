package polyjson

import (
	"encoding/json"
	"fmt"

	"github.com/bolom009/pathfind"
)

type Canvas struct {
	W float64 `json:"w"`
	H float64 `json:"h"`
}

type polygonToolJSON struct {
	Canvas   Canvas               `json:"canvas"`
	Polygons [][]pathfind.Vector2 `json:"polygons"`
}

// NewPolygonsFromJSON loads a polygon set from JSON data as exported by the
// [Polygon Constructor]: https://alaricus.github.io/PolygonConstructor/
func NewPolygonsFromJSON(jsonData []byte) (polygon []pathfind.Vector2, holes [][]pathfind.Vector2, size pathfind.Vector2, err error) {
	var jsonStruct polygonToolJSON
	err = json.Unmarshal(jsonData, &jsonStruct)
	if err != nil {
		return nil, nil, pathfind.Vector2{}, fmt.Errorf("could not unmarshal polygon JSON: %w", err)
	}

	polygon = make([]pathfind.Vector2, 0)
	for _, pp := range jsonStruct.Polygons[0] {
		polygon = append(polygon, pathfind.Vector2{X: pp.X, Y: pp.Y})
	}

	// Example hole: a small square hole in the middle: (4,4), (6,4), (6,6), (4,6)
	holes = make([][]pathfind.Vector2, 0)
	for _, pp := range jsonStruct.Polygons[1:] {
		hole := make([]pathfind.Vector2, len(pp))
		for j, pp := range pp {
			hole[j] = pathfind.Vector2{X: pp.X, Y: pp.Y}
		}

		holes = append(holes, hole)
	}

	return polygon, holes, jsonStruct.size(), nil
}

func (pj *polygonToolJSON) size() pathfind.Vector2 {
	return pathfind.Vector2{X: pj.Canvas.W, Y: pj.Canvas.H}
}
