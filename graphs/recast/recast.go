package recast

import (
	"context"
	"github.com/bolom009/delaunay"
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/graphs"
	"github.com/fzipp/astar"
	"math"
)

type Recast struct {
	polygon         []geom.Vector2
	holes           [][]geom.Vector2
	triangles       []Triangle
	costFunc        astar.CostFunc[geom.Vector2]
	visibilityGraph graphs.Graph[geom.Vector2]
}

func NewRecast(polygon []geom.Vector2, holes [][]geom.Vector2, options ...option) *Recast {
	r := &Recast{
		polygon:         polygon,
		holes:           holes,
		triangles:       make([]Triangle, 0),
		visibilityGraph: make(graphs.Graph[geom.Vector2]),
		costFunc:        heuristicEvaluation,
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (r *Recast) Generate(ctx context.Context) error {
	pp := make([]*delaunay.Point, 0)
	for _, p := range r.polygon {
		pp = append(pp, delaunay.NewPoint(p.X, p.Y))
	}

	hh := make([][]*delaunay.Point, 0)
	for _, holePoints := range r.holes {
		hole := make([]*delaunay.Point, len(holePoints))
		for i, p := range holePoints {
			hole[i] = delaunay.NewPoint(p.X, p.Y)
		}

		hh = append(hh, hole)
	}

	d := delaunay.NewSweepContext(pp, hh)

	r.triangles = convertTriangles(d.Triangulate())
	r.visibilityGraph = r.generateGraph()

	return nil
}

func (r *Recast) AggregationGraph(start, dest geom.Vector2, _ *graphs.NavOpts) graphs.Graph[geom.Vector2] {
	vis := r.visibilityGraph.Copy()

	for _, triangle := range r.triangles {
		p0 := triangle[0]
		p1 := triangle[1]
		p2 := triangle[2]

		if pointInsideTriangle(p0, p1, p2, start) {
			vis.LinkBoth(start, p0)
			vis.LinkBoth(start, p1)
			vis.LinkBoth(start, p2)
		}

		if pointInsideTriangle(p0, p1, p2, dest) {
			vis.LinkBoth(p0, dest)
			vis.LinkBoth(p1, dest)
			vis.LinkBoth(p2, dest)
		}
	}

	if isLineSegmentInsidePolygonOrHoles(r.polygon, r.holes, start, dest) {
		vis.LinkBoth(start, dest)
	}

	return vis
}

func (r *Recast) GetVisibility(_ *graphs.NavOpts) graphs.Graph[geom.Vector2] {
	return r.visibilityGraph.Copy()
}

func (r *Recast) ContainsPoint(point geom.Vector2) bool {
	return r.isInsidePolygonWithHoles(point)
}

func (r *Recast) Cost(a geom.Vector2, b geom.Vector2) float64 {
	return r.costFunc(a, b)
}

func (r *Recast) generateGraph() graphs.Graph[geom.Vector2] {
	vis := make(graphs.Graph[geom.Vector2])
	for _, triangle := range r.triangles {
		p0 := triangle[0]
		p1 := triangle[1]
		p2 := triangle[2]

		vis.LinkBoth(p0, p1)
		vis.LinkBoth(p0, p2)
		vis.LinkBoth(p1, p2)
	}

	return vis
}

func (r *Recast) isInsidePolygonWithHoles(point geom.Vector2) bool {
	if !pointInPolygon(point, r.polygon) {
		return false
	}
	for _, hole := range r.holes {
		if pointInPolygon(point, hole) {
			return false
		}
	}
	return true
}

func isLineSegmentInsidePolygonOrHoles(polygon []geom.Vector2, holes [][]geom.Vector2, lineStart, lineEnd geom.Vector2) bool {
	// Check if the line segment intersects any of the polygon edges
	if isLineSegmentInsidePolygon(polygon, lineStart, lineEnd) {
		return false
	}

	for _, hole := range holes {
		for i := 0; i < len(hole)-1; i++ {
			if doLinesIntersect(lineStart, lineEnd, hole[i], hole[i+1]) {
				return false // Intersects a hole
			}
		}
	}

	return true // Line segment is completely inside the polygon
}

// pointInPolygon checks if a point p is inside a polygon using the ray casting method.
func pointInPolygon(p geom.Vector2, poly []geom.Vector2) bool {
	inside := false
	n := len(poly)
	for i := 0; i < n; i++ {
		p1 := poly[i]
		p2 := poly[(i+1)%n]

		// Check if the edge (p1->p2) straddles the horizontal line at p.Y
		condY := (p1.Y <= p.Y && p2.Y > p.Y) || (p2.Y <= p.Y && p1.Y > p.Y)
		if condY {
			// Compute the x-coordinate of intersection of the polygon edge with the line y = p.Y
			xIntersect := p1.X + (p.Y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
			if xIntersect > p.X {
				inside = !inside
			}
		}
	}
	return inside
}

// isLineSegmentInsidePolygon function to check if a line segment is inside a polygon
func isLineSegmentInsidePolygon(polygon []geom.Vector2, lineStart, lineEnd geom.Vector2) bool {
	n := len(polygon)
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]
		if doLinesIntersect(lineStart, lineEnd, p1, p2) {
			return true // The line segment intersects the polygon edge
		}
	}

	return false
}

func doLinesIntersect(p1, p2, q1, q2 geom.Vector2) bool {
	// Convert points to segments
	// Calculate the orientation
	o1 := orientation(p1, p2, q1)
	o2 := orientation(p1, p2, q2)
	o3 := orientation(q1, q2, p1)
	o4 := orientation(q1, q2, p2)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special cases (collinear points)
	if o1 == 0 && onSegment(p1, q1, p2) {
		return true
	}
	if o2 == 0 && onSegment(p1, q2, p2) {
		return true
	}
	if o3 == 0 && onSegment(q1, p1, q2) {
		return true
	}
	if o4 == 0 && onSegment(q1, p2, q2) {
		return true
	}
	return false
}

// orientation helper function for intersection calculations
func orientation(p, q, r geom.Vector2) int {
	val := (q.Y-r.Y)*(p.X-q.X) - (q.X-r.X)*(p.Y-q.Y)
	if val == 0 {
		return 0 // Collinear
	}
	if val > 0 {
		return 1 // Clockwise
	}

	return 2 // Counterclockwise
}

// onSegment helper function for segment calculations
func onSegment(p, q, r geom.Vector2) bool {
	var (
		pX, pY = float64(p.X), float64(p.Y)
		qX, qY = float64(q.X), float64(q.Y)
		rX, rY = float64(r.X), float64(r.Y)
	)

	return qX <= math.Max(pX, rX) && qX >= math.Min(pX, rX) && qY <= math.Max(pY, rY) && qY >= math.Min(pY, rY)
}

func heuristicEvaluation(a, b geom.Vector2) float64 {
	c := a.Sub(b)
	return math.Sqrt(float64(c.X*c.X + c.Y*c.Y))
}
