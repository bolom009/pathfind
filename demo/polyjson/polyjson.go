package polyjson

import (
	"encoding/json"
	"fmt"

	"github.com/bolom009/pathfind/vec"
)

type Canvas struct {
	W float32 `json:"w"`
	H float32 `json:"h"`
}

type polygonToolJSON struct {
	Canvas   Canvas          `json:"canvas"`
	Polygons [][]vec.Vector2 `json:"polygons"`
}

// NewPolygonsFromJSON loads a polygon set from JSON data as exported by the
// [Polygon Constructor]: https://alaricus.github.io/PolygonConstructor/
func NewPolygonsFromJSON(jsonData []byte) (polygon []vec.Vector2, holes [][]vec.Vector2, size vec.Vector2, err error) {
	var jsonStruct polygonToolJSON
	err = json.Unmarshal(jsonData, &jsonStruct)
	if err != nil {
		return nil, nil, vec.Vector2{}, fmt.Errorf("could not unmarshal polygon JSON: %w", err)
	}

	polygon = make([]vec.Vector2, 0)
	for _, pp := range jsonStruct.Polygons[0] {
		polygon = append(polygon, vec.Vector2{X: pp.X, Y: pp.Y})
	}

	// Example hole: a small square hole in the middle: (4,4), (6,4), (6,6), (4,6)
	holes = make([][]vec.Vector2, 0)
	for _, pp := range jsonStruct.Polygons[1:] {
		hole := make([]vec.Vector2, len(pp))
		for j, pp := range pp {
			hole[j] = vec.Vector2{X: pp.X, Y: pp.Y}
		}

		holes = append(holes, hole)
	}

	return polygon, holes, jsonStruct.size(), nil
}

func (pj *polygonToolJSON) size() vec.Vector2 {
	return vec.Vector2{X: pj.Canvas.W, Y: pj.Canvas.H}
}
