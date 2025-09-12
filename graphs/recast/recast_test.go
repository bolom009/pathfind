package recast

import (
	"testing"

	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/demo/utils"
	"github.com/bolom009/pathfind/mesh"
)

const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":0,"y":0},{"x":120,"y":0},{"x":120,"y":340},{"x":180,"y":340},{"x":180,"y":-120},{"x":300,"y":-120},{"x":300,"y":340},{"x":360,"y":340},{"x":360,"y":0},{"x":480,"y":0},{"x":480,"y":420},{"x":300,"y":420},{"x":300,"y":540},{"x":340,"y":540},{"x":340,"y":720},{"x":140,"y":720},{"x":140,"y":540},{"x":180,"y":540},{"x":180,"y":420},{"x":0,"y":420}]]}`
const floorPlanOffset = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":696.4375,"y":67},{"x":704.4375,"y":510},{"x":121.4375,"y":509},{"x":127.4375,"y":60}],[{"x":369.4375,"y":132},{"x":368.4375,"y":209},{"x":216.4375,"y":209},{"x":216.4375,"y":135}],[{"x":378.4375,"y":342},{"x":378.4375,"y":420},{"x":250.4375,"y":421},{"x":250.4375,"y":343}],[{"x":612.4375,"y":175},{"x":643.4375,"y":251},{"x":543.4375,"y":284},{"x":487.4375,"y":222}],[{"x":633.4375,"y":370},{"x":634.4375,"y":471},{"x":472.4375,"y":452},{"x":474.4375,"y":374}]]}`

// BenchmarkAggregationGraph-16    	  823694	      1387 ns/op	    1384 B/op	       8 allocs/op
// BenchmarkAggregationGraph-16    	  801415	      1359 ns/op	    1384 B/op	       8 allocs/op
func BenchmarkAggregationGraph(b *testing.B) {
	polygon, holes, _, err := utils.NewPolygonsFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	nHoles := make([]*mesh.Hole, len(holes))
	for i, hole := range holes {
		nHoles[i] = mesh.NewHole(hole, -5.0)
	}

	var (
		recastGraph = NewRecast(mesh.NewPolygon(polygon, 35), nHoles, WithSearchOutOfArea(true))
		// out area
		start = geom.Vector2{X: 60, Y: 10}
		dest  = geom.Vector2{X: 425, Y: 10}
		//inside area
		//start = geom.Vector2{X: 25, Y: 25}
		//dest  = geom.Vector2{X: 430, Y: 200}
		// direct
		//start = geom.Vector2{X: 60, Y: 60}
		//dest  = geom.Vector2{X: 60, Y: 100}
	)

	_ = recastGraph.Generate(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = recastGraph.AggregationGraph(start, dest, nil)
	}
}

func BenchmarkAggregationGraphStartOffset(b *testing.B) {
	polygon, holes, _, err := utils.NewPolygonsFromJSON([]byte(floorPlanOffset))
	if err != nil {
		panic(err)
	}

	nHoles := make([]*mesh.Hole, len(holes))
	for i, hole := range holes {
		nHoles[i] = mesh.NewHole(hole, -20.0)
	}

	var (
		recastGraph = NewRecast(mesh.NewPolygon(polygon, 20), nHoles, WithSearchOutOfArea(true))
		start       = geom.Vector2{X: 236, Y: 130}
		dest        = geom.Vector2{X: 656, Y: 446}
	)

	_ = recastGraph.Generate(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = recastGraph.AggregationGraph(start, dest, nil)
	}
}

func TestAggregationGraph(t *testing.T) {
	polygon, holes, _, err := utils.NewPolygonsFromJSON([]byte(floorPlan))
	if err != nil {
		panic(err)
	}

	nHoles := make([]*mesh.Hole, len(holes))
	for i, hole := range holes {
		nHoles[i] = mesh.NewHole(hole, -5.0)
	}

	var (
		recastGraph = NewRecast(mesh.NewPolygon(polygon, 35), nHoles, WithSearchOutOfArea(true))
		start       = geom.Vector2{X: 60, Y: 10}
		dest        = geom.Vector2{X: 425, Y: 10}
	)

	_ = recastGraph.Generate(nil)

	_ = recastGraph.AggregationGraph(start, dest, nil)
	//for i, edges := range vis {
	//	fmt.Println(i, cap(edges))
	//}
}
