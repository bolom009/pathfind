package grid

import (
	"context"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/vec"
	"github.com/fzipp/astar"
	"math"
)

type Grid struct {
	polygon         []vec.Vector2
	holes           [][]vec.Vector2
	squares         []Square
	visibilityGraph graphs.Graph[vec.Vector2]
	squareSize      float32
	costFunc        astar.CostFunc[vec.Vector2]
}

func NewGrid(polygon []vec.Vector2, holes [][]vec.Vector2, squareSize float32, options ...option) *Grid {
	g := &Grid{
		polygon:         polygon,
		holes:           holes,
		squareSize:      squareSize,
		squares:         make([]Square, 0),
		visibilityGraph: make(graphs.Graph[vec.Vector2]),
		costFunc:        heuristicEvaluation,
	}

	for _, option := range options {
		option(g)
	}

	return g
}

func (g *Grid) Generate(_ context.Context) error {
	g.squares = g.generateSquares()
	g.visibilityGraph = g.generateGraph()

	return nil
}

func (g *Grid) ContainsPoint(point vec.Vector2) bool {
	return g.isInsidePolygonWithHoles(point)
}

func (g *Grid) GetVisibility() graphs.Graph[vec.Vector2] {
	return g.visibilityGraph.Copy()
}

func (g *Grid) Cost(a, b vec.Vector2) float64 { return g.costFunc(a, b) }

// AggregationGraph add start and dest points to existing pathfinder graph
func (g *Grid) AggregationGraph(start, dest vec.Vector2) graphs.Graph[vec.Vector2] {
	vis := g.visibilityGraph.Copy()

	for _, square := range g.squares {
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
			g.addPairs(vis, start, []gridPair{{a, isA}, {b, isB}, {c, isC}, {d, isD}})

			if isA && g.isLineSegmentInsidePolygon(a, start) {
				vis.Link(a, start)
			}
			if isB && g.isLineSegmentInsidePolygon(b, start) {
				vis.Link(b, start)
			}
			if isC && g.isLineSegmentInsidePolygon(c, start) {
				vis.Link(c, start)
			}
			if isD && g.isLineSegmentInsidePolygon(d, start) {
				vis.Link(d, start)
			}
		}

		if square.isPointInsideSquare(dest) {
			g.addPairs(vis, dest, []gridPair{{a, isA}, {b, isB}, {c, isC}, {d, isD}})

			if isA && g.isLineSegmentInsidePolygon(a, dest) {
				vis.Link(a, dest)
			}
			if isB && g.isLineSegmentInsidePolygon(b, dest) {
				vis.Link(b, dest)
			}
			if isC && g.isLineSegmentInsidePolygon(c, dest) {
				vis.Link(c, dest)
			}
			if isD && g.isLineSegmentInsidePolygon(d, dest) {
				vis.Link(d, dest)
			}
		}
	}

	return vis
}

func (g *Grid) Squares() []Square {
	cSquares := make([]Square, len(g.squares))
	for i, square := range g.squares {
		cSquares[i] = square
	}

	return cSquares
}

// generateGraph create visibility graph based on squares
func (g *Grid) generateGraph() graphs.Graph[vec.Vector2] {
	vis := make(graphs.Graph[vec.Vector2])
	for _, square := range g.squares {
		var (
			a, b, c, d         = square.A, square.B, square.C, square.D
			isA, isB, isC, isD = square.isA, square.isB, square.isC, square.isD
		)

		if isA {
			g.addPairs(vis, a, []gridPair{{b, isB}, {d, isD}})
		}
		if isB {
			g.addPairs(vis, b, []gridPair{{a, isA}, {c, isC}})
		}
		if isC {
			g.addPairs(vis, c, []gridPair{{b, isB}, {d, isD}})
		}
		if isD {
			g.addPairs(vis, d, []gridPair{{a, isA}, {c, isC}})
		}
		if isA && isC && g.isLineSegmentInsidePolygon(a, c) {
			vis.Link(a, c)
			vis.Link(c, a)
		}
		if isB && isD && g.isLineSegmentInsidePolygon(b, d) {
			vis.Link(b, d)
			vis.Link(d, b)
		}
	}

	c := 0
	for _, p := range vis {
		c += len(p)
	}

	return vis
}

// generateSquares create list of squares based on polygon, holes and squareSize
// each square has info about vertices and if each vertex inside polygon to calculate graph
func (g *Grid) generateSquares() []Square {
	var (
		squareSize = g.squareSize
		polygon    = g.polygon
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

				center = vec.Vector2{X: x + squareSize/2.0, Y: y + squareSize/2.0}
				a      = vec.Vector2{X: x, Y: y}
				b      = vec.Vector2{X: x + squareSize, Y: y}
				c      = vec.Vector2{X: x + squareSize, Y: y + squareSize}
				d      = vec.Vector2{X: x, Y: y + squareSize}
			)

			if g.isInsidePolygonWithHoles(a) {
				isA = true
			}

			if g.isInsidePolygonWithHoles(b) {
				isB = true
			}

			if g.isInsidePolygonWithHoles(c) {
				isC = true
			}

			if g.isInsidePolygonWithHoles(d) {
				isD = true
			}

			if !isA && !isB && !isC && !isD {
				continue
			}

			if g.isInsidePolygonWithHoles(center) {
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

// pointInPolygon checks if a point p is inside a polygon using the ray casting method.
func pointInPolygon(p vec.Vector2, poly []vec.Vector2) bool {
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
func (g *Grid) isInsidePolygonWithHoles(point vec.Vector2) bool {
	if !pointInPolygon(point, g.polygon) {
		return false
	}
	for _, hole := range g.holes {
		if pointInPolygon(point, hole) {
			return false
		}
	}
	return true
}

// isLineSegmentInsidePolygon function to check if a line segment is inside a polygon
func (g *Grid) isLineSegmentInsidePolygon(lineStart, lineEnd vec.Vector2) bool {
	var (
		polygon = g.polygon
		holes   = g.holes
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
func doLinesIntersect(p1, p2, q1, q2 vec.Vector2) bool {
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
func orientation(p, q, r vec.Vector2) int {
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
func onSegment(p, q, r vec.Vector2) bool {
	var (
		pX, pY = float64(p.X), float64(p.Y)
		qX, qY = float64(q.X), float64(q.Y)
		rX, rY = float64(r.X), float64(r.Y)
	)

	return qX <= math.Max(pX, rX) && qX >= math.Min(pX, rX) && qY <= math.Max(pY, rY) && qY >= math.Min(pY, rY)
}

func heuristicEvaluation(a, b vec.Vector2) float64 {
	c := a.Sub(b)
	return math.Sqrt(float64(c.X*c.X + c.Y*c.Y))
}
