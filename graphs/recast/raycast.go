package recast

import (
	"math"

	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/mesh"
)

type Hole struct {
	Points   []geom.Vector2
	Viewable bool
}

type Line struct {
	A geom.Vector2
	B geom.Vector2
}

func (line *Line) Center() geom.Vector2 {
	return geom.Vector2{
		X: (line.A.X + line.B.X) / 2,
		Y: (line.A.Y + line.B.Y) / 2,
	}
}

func (line *Line) Length() float32 {
	dx := line.B.X - line.A.X
	dy := line.B.Y - line.A.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// BoundingBox represents an axis-aligned bounding box
type BoundingBox struct {
	MinX, MinY, MaxX, MaxY float32
}

type Raycast struct {
	polygon   []geom.Vector2
	holes     []*mesh.Hole
	outerBBox BoundingBox
	holesBBox []BoundingBox
}

func NewRaycast(polygon []geom.Vector2, holes []*mesh.Hole) *Raycast {
	obj := &Raycast{
		polygon:   polygon,
		holes:     holes,
		outerBBox: getBoundingBox(polygon),
		holesBBox: make([]BoundingBox, len(holes)),
	}

	for i, hole := range holes {
		obj.holesBBox[i] = getBoundingBox(hole.Points())
	}

	return obj
}

// CheckLineIntersectsPolygon checks if a line intersects with the polygon or its holes
func (v *Raycast) CheckLineIntersectsPolygon(lineStart, lineEnd geom.Vector2) bool {
	// Bounding box check for outer polygon
	if !lineIntersectsBoundingBox(lineStart, lineEnd, v.outerBBox) {
		return false
	}

	// Check intersection with outer polygon
	if lineIntersectsPolygonOutline(lineStart, lineEnd, v.polygon) {
		return true
	}

	for i, hole := range v.holes {
		if hole.Viewable() {
			continue
		}

		if !lineIntersectsBoundingBox(lineStart, lineEnd, v.holesBBox[i]) {
			continue
		}
		if lineIntersectsPolygonOutline(lineStart, lineEnd, hole.Points()) {
			return true
		}
	}

	return false
}

// lineIntersectsBoundingBox checks if line intersects with bounding box
func lineIntersectsBoundingBox(p1, p2 geom.Vector2, box BoundingBox) bool {
	// Check if either point is inside box
	if pointInBoundingBox(p1, box) || pointInBoundingBox(p2, box) {
		return true
	}
	// Check if line intersects any of the box's edges
	boxPoints := []geom.Vector2{
		{box.MinX, box.MinY},
		{box.MaxX, box.MinY},
		{box.MaxX, box.MaxY},
		{box.MinX, box.MaxY},
	}
	edges := [][]geom.Vector2{
		{boxPoints[0], boxPoints[1]},
		{boxPoints[1], boxPoints[2]},
		{boxPoints[2], boxPoints[3]},
		{boxPoints[3], boxPoints[0]},
	}
	for _, edge := range edges {
		if lineSegmentIntersection(p1, p2, edge[0], edge[1]) {
			return true
		}
	}
	return false
}

// lineIntersectsPolygonOutline checks line intersection with polygon outline
func lineIntersectsPolygonOutline(lineStart, lineEnd geom.Vector2, polygonPoints []geom.Vector2) bool {
	n := len(polygonPoints)
	for i := 0; i < n; i++ {
		p1 := polygonPoints[i]
		p2 := polygonPoints[(i+1)%n] // Wrap around to the first point
		if lineSegmentIntersection(lineStart, lineEnd, p1, p2) {
			return true
		}
	}
	return false
}

// lineSegmentIntersection checks if two line segments intersect
func lineSegmentIntersection(p1, p2, p3, p4 geom.Vector2) bool {
	s1 := p4.Y - p3.Y
	s2 := p2.X - p1.X
	s3 := p4.X - p3.X
	s4 := p2.Y - p1.Y

	denom := s1*s2 - s3*s4
	if denom == 0 {
		return false // Parallel lines
	}

	s5 := p1.Y - p3.Y
	s6 := p1.X - p3.X

	ua := (s3*s5 - s1*s6) / denom
	ub := (s2*s5 - s4*s6) / denom
	return ua >= 0 && ua <= 1 && ub >= 0 && ub <= 1
}

// pointInBoundingBox checks if a point is inside a bounding box
func pointInBoundingBox(p geom.Vector2, b BoundingBox) bool {
	return p.X >= b.MinX && p.X <= b.MaxX && p.Y >= b.MinY && p.Y <= b.MaxY
}

// getBoundingBox calculates the bounding box for a polygon
func getBoundingBox(points []geom.Vector2) BoundingBox {
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := points[0].X, points[0].Y

	for _, p := range points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return BoundingBox{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
}
