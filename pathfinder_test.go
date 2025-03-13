package pathfind

import (
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/demo/utils"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/graphs/recast"
	"github.com/bolom009/pathfind/mesh"
	"testing"
)

const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":0,"y":0},{"x":120,"y":0},{"x":120,"y":340},{"x":180,"y":340},{"x":180,"y":-120},{"x":300,"y":-120},{"x":300,"y":340},{"x":360,"y":340},{"x":360,"y":0},{"x":480,"y":0},{"x":480,"y":420},{"x":300,"y":420},{"x":300,"y":540},{"x":340,"y":540},{"x":340,"y":720},{"x":140,"y":720},{"x":140,"y":540},{"x":180,"y":540},{"x":180,"y":420},{"x":0,"y":420}]]}`

func BenchmarkPathfinder_Path(b *testing.B) {
	polygon, holes, _, err := utils.NewPolygonsFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	nHoles := make([]*mesh.Hole, len(holes))
	for i, hole := range holes {
		nHoles[i] = mesh.NewHole(hole, -5.0)
	}

	var (
		start       = geom.Vector2{X: 60, Y: 10}
		dest        = geom.Vector2{X: 425, Y: 10}
		recastGraph = recast.NewRecast(mesh.NewPolygon(polygon, 35), nHoles, recast.WithSearchOutOfArea(true))
		pathfinder  = NewPathfinder[geom.Vector2]([]graphs.NavGraph[geom.Vector2]{
			recastGraph,
		})
	)

	_ = pathfinder.Initialize(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = pathfinder.Path(0, start, dest)
	}
}
