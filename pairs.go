package pathfind

type pair struct {
	point Vector2
	is    bool
}

func (p *Pathfinder) addPairs(vis graph[Vector2], dest Vector2, pairs []pair) {
	points := make([]Vector2, 0)
	for _, pair := range pairs {
		if pair.is && p.isLineSegmentInsidePolygon(dest, pair.point) {
			points = append(points, pair.point)
		}
	}

	// add visible nodes
	for _, point := range points {
		vis.link(dest, point)
	}
}
