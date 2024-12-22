package pathfind

import (
	"fmt"
	"github.com/fzipp/astar"
	"math"
)

// Pathfinder represent struct to operate with pathfinding
// Graph of current pathfinder represent as list of squares
type Pathfinder struct {
	polygon         []Vector2
	holes           [][]Vector2
	squares         []Square
	squareSize      float64
	visibilityGraph graph[Vector2]
	constFunc       astar.CostFunc[Vector2]
}

// NewPathfinder constructor to create pathfinder struct
// polygons contains two parts: polygon and holes
func NewPathfinder(polygon []Vector2, holes [][]Vector2, squareSize float64, options ...Option) *Pathfinder {
	p := &Pathfinder{
		polygon:    polygon,
		holes:      holes,
		squareSize: squareSize,
		constFunc:  heuristicEvaluation,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// Initialize split polygon with holes to square and prepare visibility graph
func (p *Pathfinder) Initialize() {
	p.squares = p.generateSquares()
	p.visibilityGraph = p.generateGraph()
}

// Path finds the shortest path from start to dest
// If dest is outside the polygon set it will be clamped to the nearest polygon point.
// The function returns nil if no path exists
func (p *Pathfinder) Path(start, dest Vector2) []Vector2 {
	if !p.isInsidePolygonWithHoles(dest) {
		return nil
	}

	vis := p.addStartDestToGraph(start, dest)

	return astar.FindPath[Vector2](vis, start, dest, p.constFunc, p.constFunc)
}

// Squares return list of generated squares
func (p *Pathfinder) Squares() []Square {
	return p.squares
}

// Graph return generated graph based on square list
func (p *Pathfinder) Graph() graph[Vector2] {
	return p.visibilityGraph.copy()
}

// GraphWithSearchPath return generated graph with path nodes
func (p *Pathfinder) GraphWithSearchPath(start, dest Vector2) graph[Vector2] {
	return p.addStartDestToGraph(start, dest)
}

// generateGraph create visibility graph based on squares
func (p *Pathfinder) generateGraph() graph[Vector2] {
	vis := make(graph[Vector2])
	for _, square := range p.squares {
		var (
			a, b, c, d         = square.A, square.B, square.C, square.D
			isA, isB, isC, isD = square.isA, square.isB, square.isC, square.isD
		)

		if isA {
			p.addPairs(vis, a, []pair{{b, isB}, {d, isD}})
		}
		if isB {
			p.addPairs(vis, b, []pair{{a, isA}, {c, isC}})
		}
		if isC {
			p.addPairs(vis, c, []pair{{b, isB}, {d, isD}})
		}
		if isD {
			p.addPairs(vis, d, []pair{{a, isA}, {c, isC}})
		}
		if isA && isC && p.isLineSegmentInsidePolygon(a, c) {
			vis.link(a, c)
			vis.link(c, a)
		}
		if isB && isD && p.isLineSegmentInsidePolygon(b, d) {
			vis.link(b, d)
			vis.link(d, b)
		}
	}

	c := 0
	for _, p := range vis {
		c += len(p)
	}

	fmt.Println("==> GRAPH", c)

	return vis
}

// generateSquares create list of squares based on polygon, holes and squareSize
// each square has info about vertices and if each vertex inside polygon to calculate graph
func (p *Pathfinder) generateSquares() []Square {
	var (
		squareSize = p.squareSize
		polygon    = p.polygon
	)

	// Compute bounding box of the outer polygon
	minX, maxX := polygon[0].X, polygon[0].X
	minY, maxY := polygon[0].Y, polygon[0].Y
	for _, pt := range polygon {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	squares := make([]Square, 0, int(squareSize*squareSize))

	// Iterate over the grid
	for x := minX; x < maxX; x += squareSize {
		for y := minY; y < maxY; y += squareSize {
			var (
				isA, isB, isC, isD, isCenter bool

				center = Vector2{X: x + squareSize/2.0, Y: y + squareSize/2.0}
				a      = Vector2{X: x, Y: y}
				b      = Vector2{X: x + squareSize, Y: y}
				c      = Vector2{X: x + squareSize, Y: y + squareSize}
				d      = Vector2{X: x, Y: y + squareSize}
			)

			if p.isInsidePolygonWithHoles(a) {
				isA = true
			}

			if p.isInsidePolygonWithHoles(b) {
				isB = true
			}

			if p.isInsidePolygonWithHoles(c) {
				isC = true
			}

			if p.isInsidePolygonWithHoles(d) {
				isD = true
			}

			if !isA && !isB && !isC && !isD {
				continue
			}

			if p.isInsidePolygonWithHoles(center) {
				isCenter = true
			}

			if isA || isB || isC || isD || isCenter {
				// The square’s corners
				square := Square{
					A:        a,
					B:        b,
					C:        c,
					D:        d,
					isA:      isA,
					isB:      isB,
					isC:      isC,
					isD:      isD,
					isCenter: isCenter,
					Center:   center,
				}

				squares = append(squares, square)
			}
		}
	}

	return squares
}

// addStartDestToGraph add start and dest points to existing pathfinder graph
func (p *Pathfinder) addStartDestToGraph(start, dest Vector2) graph[Vector2] {
	vis := p.visibilityGraph.copy()

	for _, square := range p.squares {
		var (
			a   = square.A
			b   = square.B
			c   = square.C
			d   = square.D
			isA = square.isA
			isB = square.isB
			isC = square.isC
			isD = square.isD
		)

		if square.isPointInsideSquare(start) {
			p.addPairs(vis, start, []pair{{a, isA}, {b, isB}, {c, isC}, {d, isD}})

			if isA && p.isLineSegmentInsidePolygon(a, start) {
				vis.link(a, start)
			}
			if isB && p.isLineSegmentInsidePolygon(b, start) {
				vis.link(b, start)
			}
			if isC && p.isLineSegmentInsidePolygon(c, start) {
				vis.link(c, start)
			}
			if isD && p.isLineSegmentInsidePolygon(d, start) {
				vis.link(d, start)
			}
		}

		if square.isPointInsideSquare(dest) {
			p.addPairs(vis, dest, []pair{{a, isA}, {b, isB}, {c, isC}, {d, isD}})

			if isA && p.isLineSegmentInsidePolygon(a, dest) {
				vis.link(a, dest)
			}
			if isB && p.isLineSegmentInsidePolygon(b, dest) {
				vis.link(b, dest)
			}
			if isC && p.isLineSegmentInsidePolygon(c, dest) {
				vis.link(c, dest)
			}
			if isD && p.isLineSegmentInsidePolygon(d, dest) {
				vis.link(d, dest)
			}
		}
	}

	return vis
}

// pointInPolygon checks if a point p is inside a polygon using the ray casting method.
func pointInPolygon(p Vector2, poly []Vector2) bool {
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

// isInsidePolygonWithHoles checks if p is inside the outer polygon but not inside any holes
func (p *Pathfinder) isInsidePolygonWithHoles(point Vector2) bool {
	if !pointInPolygon(point, p.polygon) {
		return false
	}
	for _, hole := range p.holes {
		if pointInPolygon(point, hole) {
			return false
		}
	}
	return true
}

// isLineSegmentInsidePolygon function to check if a line segment is inside a polygon
func (p *Pathfinder) isLineSegmentInsidePolygon(lineStart, lineEnd Vector2) bool {
	var (
		polygon = p.polygon
		holes   = p.holes
	)

	// Check if the line segment intersects any of the polygon edges
	n := len(polygon)
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]
		if doLinesIntersect(lineStart, lineEnd, p1, p2) {
			return false // The line segment intersects the polygon edge
		}
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

// Function to check if two lines intersect
func doLinesIntersect(p1, p2, q1, q2 Vector2) bool {
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
func orientation(p, q, r Vector2) int {
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
func onSegment(p, q, r Vector2) bool {
	return q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) && q.Y <= math.Max(p.Y, r.Y) && q.Y >= math.Min(p.Y, r.Y)
}

func heuristicEvaluation(a, b Vector2) float64 {
	c := a.Sub(b)
	return math.Sqrt(c.X*c.X + c.Y*c.Y)
}
