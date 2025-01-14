package grid

import (
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/vec"
)

type gridPair struct {
	point vec.Vector2
	is    bool
}

func (g *Grid) addPairs(vis graphs.Graph[vec.Vector2], dest vec.Vector2, pairs []gridPair) {
	points := make([]vec.Vector2, 0)
	for _, pair := range pairs {
		if pair.is && g.isLineSegmentInsidePolygonOrHoles(dest, pair.point) {
			points = append(points, pair.point)
		}
	}

	// add visible nodes
	for _, point := range points {
		vis.Link(dest, point)
	}
}
