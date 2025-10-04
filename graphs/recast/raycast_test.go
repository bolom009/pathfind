package recast

import (
	"testing"

	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/demo/utils"
	"github.com/bolom009/pathfind/mesh"
)

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
