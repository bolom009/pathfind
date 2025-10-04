package recast

import (
	"math"

	"github.com/bolom009/geom"
)

const eps = 1e-9

// isPolygonAFullyInsideB returns true if polygonA is fully inside polygonB.
// If allowTouchOnBoundary==true then points on B's edges count as inside.
// If false, any touching or edge intersection returns false.
func isPolygonAFullyInsideB(polygonA, polygonB []geom.Vector2, allowTouchOnBoundary bool) bool {
	if !hasEnoughPoints(polygonA) || !hasEnoughPoints(polygonB) {
		return false
	}

	// Quick bbox reject
	if !boundsContains(getBounds(polygonB), getBounds(polygonA)) {
		return false
	}

	// 1) Every vertex of A must be inside B (or on boundary if allowed)
	for _, v := range polygonA {
		pip := pointInPolygonF(v, polygonB)
		if pip == 0 {
			return false // outside
		}
		if !allowTouchOnBoundary && pip == -1 {
			return false // on edge not allowed
		}
	}

	// 2) No edge of A should properly intersect any edge of B.
	if edgesIntersectProperly(polygonA, polygonB, allowTouchOnBoundary) {
		return false
	}

	return true
}

func hasEnoughPoints(p []geom.Vector2) bool { return len(p) >= 3 }

type rect struct {
	minX, minY, maxX, maxY float32
}

func getBounds(p []geom.Vector2) rect {
	minX := float32(math.Inf(1))
	minY := float32(math.Inf(1))
	maxX := float32(math.Inf(-1))
	maxY := float32(math.Inf(-1))
	for _, v := range p {
		if v.X < minX {
			minX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}
	return rect{minX, minY, maxX, maxY}
}

func boundsContains(outer, inner rect) bool {
	return !(inner.minX < outer.minX-eps || inner.maxX > outer.maxX+eps ||
		inner.minY < outer.minY-eps || inner.maxY > outer.maxY+eps)
}

// pointInPolygonF: ray-casting (even-odd) with on-edge detection.
// returns: 1 = inside, 0 = outside, -1 = on-edge
func pointInPolygonF(pt geom.Vector2, poly []geom.Vector2) int {
	inside := false
	n := len(poly)
	if n == 0 {
		return 0
	}
	j := n - 1
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[j]

		if isPointOnSegment(pt, a, b) {
			return -1
		}

		intersect := ((a.Y > pt.Y) != (b.Y > pt.Y)) &&
			(pt.X < (b.X-a.X)*(pt.Y-a.Y)/(b.Y-a.Y)+a.X)
		if intersect {
			inside = !inside
		}
		j = i
	}
	if inside {
		return 1
	}
	return 0
}

func isPointOnSegment(p, a, b geom.Vector2) bool {
	// bounding box test
	if float64(p.X) < math.Min(float64(a.X), float64(b.X))-eps || float64(p.X) > math.Max(float64(a.X), float64(b.X))+eps ||
		float64(p.Y) < math.Min(float64(a.Y), float64(b.Y))-eps || float64(p.Y) > math.Max(float64(a.Y), float64(b.Y))+eps {
		return false
	}
	cross := (p.X-a.X)*(b.Y-a.Y) - (p.Y-a.Y)*(b.X-a.X)
	if math.Abs(float64(cross)) > eps {
		return false
	}
	return true
}

// returns true if there is any "proper" intersection between edges of A and B.
// If allowTouchOnBoundary==true, then touching at a point or overlapping collinear segments that lie on B's boundary are allowed.
func edgesIntersectProperly(A, B []geom.Vector2, allowTouchOnBoundary bool) bool {
	na := len(A)
	nb := len(B)
	if na == 0 || nb == 0 {
		return false
	}
	for i := 0; i < na; i++ {
		a1 := A[i]
		a2 := A[(i+1)%na]
		for j := 0; j < nb; j++ {
			b1 := B[j]
			b2 := B[(j+1)%nb]
			if segmentsProperlyIntersect(a1, a2, b1, b2, allowTouchOnBoundary) {
				return true
			}
		}
	}
	return false
}

// Checks if two segments intersect in a way that violates "A inside B".
// If allowTouchOnBoundary==true, intersections that are only touches or collinear overlaps on boundary are allowed.
func segmentsProperlyIntersect(p1, p2, p3, p4 geom.Vector2, allowTouchOnBoundary bool) bool {
	d1 := float64(cross(p3, p4, p1))
	d2 := float64(cross(p3, p4, p2))
	d3 := float64(cross(p1, p2, p3))
	d4 := float64(cross(p1, p2, p4))

	// Proper intersection where segments cross transversely
	if ((d1 > eps && d2 < -eps) || (d1 < -eps && d2 > eps)) &&
		((d3 > eps && d4 < -eps) || (d3 < -eps && d4 > eps)) {
		return true
	}

	// Collinear or endpoint cases
	if math.Abs(d1) <= eps && isPointOnSegment(p1, p3, p4) {
		if !allowTouchOnBoundary {
			return true
		}
		return false
	}
	if math.Abs(d2) <= eps && isPointOnSegment(p2, p3, p4) {
		if !allowTouchOnBoundary {
			return true
		}
		return false
	}
	if math.Abs(d3) <= eps && isPointOnSegment(p3, p1, p2) {
		if !allowTouchOnBoundary {
			return true
		}
		return false
	}
	if math.Abs(d4) <= eps && isPointOnSegment(p4, p1, p2) {
		if !allowTouchOnBoundary {
			return true
		}
		return false
	}

	return false
}

func cross(a, b, c geom.Vector2) float32 {
	return (c.X-a.X)*(b.Y-a.Y) - (c.Y-a.Y)*(b.X-a.X)
}
