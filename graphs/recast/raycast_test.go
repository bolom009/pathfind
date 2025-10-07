package recast

import (
	"testing"

	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/demo/utils"
	"github.com/bolom009/pathfind/mesh"
)

const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":0,"y":0},{"x":120,"y":0},{"x":120,"y":340},{"x":180,"y":340},{"x":180,"y":-120},{"x":300,"y":-120},{"x":300,"y":340},{"x":360,"y":340},{"x":360,"y":0},{"x":480,"y":0},{"x":480,"y":420},{"x":300,"y":420},{"x":300,"y":540},{"x":340,"y":540},{"x":340,"y":720},{"x":140,"y":720},{"x":140,"y":540},{"x":180,"y":540},{"x":180,"y":420},{"x":0,"y":420}]]}`

func BenchmarkCheckLineIntersectsPolygon(b *testing.B) {
	polygon, holes, _, err := utils.NewPolygonsFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	newHoles := make([]*mesh.Hole, len(holes))
	for i, points := range holes {
		newHoles[i] = mesh.NewObstacle(points, 0, i%2 == 0)
	}

	var (
		//start   = geom.Vector2{172, 104}
		//end     = geom.Vector2{178, 452}
		//start = geom.Vector2{286, 116}
		//end   = geom.Vector2{284, 228}

		start = geom.Vector2{176, 76}
		end   = geom.Vector2{656, 482}

		raycast = NewRaycast(polygon, newHoles)
	)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = raycast.CheckLineIntersectsPolygon(start, end)
	}
}
