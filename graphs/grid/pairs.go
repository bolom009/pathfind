package grid

import (
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/graphs"
)

type gridPair struct {
	point geom.Vector2
	is    bool
}

func (g *Grid) addPairs(vis graphs.Graph[geom.Vector2], dest geom.Vector2, pairs []gridPair) {
	points := make([]geom.Vector2, 0)
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
