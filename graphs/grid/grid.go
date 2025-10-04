package grid

import (
	"context"
	"math"

	"github.com/bolom009/astar"
	"github.com/bolom009/geom"
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/obstacles"
)

type Grid struct {
	polygon         []geom.Vector2
	holes           [][]geom.Vector2
	squares         []Square
	visSquares      []Square
	visibilityGraph graphs.Graph[geom.Vector2]
	squareSize      float32
	costFunc        astar.CostFunc[geom.Vector2]
	offset          geom.Vector2
}

func NewGrid(polygon []geom.Vector2, holes [][]geom.Vector2, squareSize float32, options ...option) *Grid {
	g := &Grid{
		polygon:         polygon,
		holes:           holes,
		squareSize:      squareSize,
		squares:         make([]Square, 0),
		visibilityGraph: make(graphs.Graph[geom.Vector2]),
		costFunc:        heuristicEvaluation,
	}

	for _, option := range options {
		option(g)
	}

	return g
}

func (g *Grid) Generate(_ context.Context) error {
	g.squares, g.visSquares = g.generateSquares()
	g.visibilityGraph = g.generateGraph()

	return nil
}

func (g *Grid) ContainsPoint(point geom.Vector2) bool {
	return g.isInsidePolygonWithHoles(point)
}

func (g *Grid) GetVisibility(navOpts *graphs.NavOpts) graphs.Graph[geom.Vector2] {
	vis := g.visibilityGraph.Copy()
	if navOpts == nil {
		return vis
	}

	if navOpts.Obstacles != nil {
		// cut graph with obstacles
		g.updateGraphWithObstacles(vis, navOpts.Obstacles)
	}

	return vis
}

func (g *Grid) Cost(a, b geom.Vector2) float32 {
	return g.costFunc(a, b)
}

func (g *Grid) IsRaycastHit(start, end geom.Vector2) bool {
	return false
}

func (g *Grid) GetClosestPoint(point geom.Vector2) (geom.Vector2, bool) {
	closest := float32(math.MaxFloat32)
	closestPoint := geom.Vector2{}
	for _, visSq := range g.visSquares {
		for _, v := range []geom.Vector2{visSq.A, visSq.B, visSq.C, visSq.D, visSq.Center} {
			dist := geom.Distance(point, v)
			if dist < closest {
				closest = dist
				closestPoint = v
			}
		}
	}

	if closest == math.MaxFloat32 {
		return closestPoint, false
	}

	return closestPoint, true
}

// AggregationGraph add start and dest points to existing pathfinder graph
func (g *Grid) AggregationGraph(start, dest geom.Vector2, navOpts *graphs.NavOpts) graphs.Graph[geom.Vector2] {
	vis := g.visibilityGraph.Copy()

	// add start & dest points to graph
	g.addStartDestPointsToGraph(vis, start, dest)

	if navOpts != nil {
		if navOpts.Obstacles != nil {
			// cut graph with obstacles
			g.updateGraphWithObstacles(vis, navOpts.Obstacles)
		}
	}

	return vis
}

// Squares return copied list of squares
func (g *Grid) Squares() []Square {
	cSquares := make([]Square, len(g.squares))
	copy(cSquares, g.squares)

	return cSquares
}

// VisibleSquares return copied list of squares
func (g *Grid) VisibleSquares() []Square {
	cSquares := make([]Square, len(g.visSquares))
	copy(cSquares, g.visSquares)

	return cSquares
}

func (g *Grid) addStartDestPointsToGraph(vis graphs.Graph[geom.Vector2], start geom.Vector2, dest geom.Vector2) {
	for _, square := range g.visSquares {
		var (
			a = square.A
			b = square.B
			c = square.C
			d = square.D
		)

		if square.isPointInsideSquare(start) {
			vis.LinkBoth(a, start)
			vis.LinkBoth(b, start)
			vis.LinkBoth(c, start)
			vis.LinkBoth(d, start)
		}

		if square.isPointInsideSquare(dest) {
			vis.LinkBoth(a, dest)
			vis.LinkBoth(b, dest)
			vis.LinkBoth(c, dest)
			vis.LinkBoth(d, dest)
		}
	}
}

func (g *Grid) updateGraphWithObstacles(vis graphs.Graph[geom.Vector2], obstacles []obstacles.Obstacle) {
	for _, square := range g.visSquares {
		for _, obstacle := range obstacles {
			var (
				a               = square.A
				b               = square.B
				c               = square.C
				d               = square.D
				obstaclePolygon = obstacle.GetPolygon()
			)

			// is squire center around or inside obstacle
			if !obstacle.IsPointAround(square.Center, g.squareSize) {
				continue
			}

			// check edges list
			for _, neighbour := range vis.Neighbours(a) {
				if isLineSegmentInsidePolygon(obstaclePolygon, a, neighbour) {
					vis.DeleteNeighbour(a, neighbour)
				}
			}

			for _, neighbour := range vis.Neighbours(b) {
				if isLineSegmentInsidePolygon(obstaclePolygon, b, neighbour) {
					vis.DeleteNeighbour(b, neighbour)
				}
			}

			for _, neighbour := range vis.Neighbours(c) {
				if isLineSegmentInsidePolygon(obstaclePolygon, c, neighbour) {
					vis.DeleteNeighbour(c, neighbour)
				}
			}

			for _, neighbour := range vis.Neighbours(d) {
				if isLineSegmentInsidePolygon(obstaclePolygon, d, neighbour) {
					vis.DeleteNeighbour(d, neighbour)
				}
			}

			// check vertex list
			if pointInPolygon(a, obstaclePolygon) {
				vis.DeleteNode(a)
			}
			if pointInPolygon(b, obstaclePolygon) {
				vis.DeleteNode(b)
			}
			if pointInPolygon(c, obstaclePolygon) {
				vis.DeleteNode(c)
			}
			if pointInPolygon(d, obstaclePolygon) {
				vis.DeleteNode(d)
			}
		}
	}
}

// generateGraph create visibility graph based on squares
func (g *Grid) generateGraph() graphs.Graph[geom.Vector2] {
	vis := make(graphs.Graph[geom.Vector2])
	for _, square := range g.visSquares {
		var (
			a, b, c, d = square.A, square.B, square.C, square.D
		)

		vis.LinkBoth(a, b)
		vis.LinkBoth(a, d)
		vis.LinkBoth(c, b)
		vis.LinkBoth(c, d)
		vis.LinkBoth(a, c)
		vis.LinkBoth(b, d)
	}

	return vis
}

// generateSquares create list of squares based on polygon, holes and squareSize
// each square has info about vertices and if each vertex inside polygon to calculate graph
func (g *Grid) generateSquares() ([]Square, []Square) {
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
	visSquares := make([]Square, 0, int(squareSize*squareSize))

	// Iterate over the grid
	for x := minX; x < maxX; x += squareSize {
		for y := minY; y < maxY; y += squareSize {
			var (
				isA, isB, isC, isD, isCenter bool

				center = geom.Vector2{X: x + g.offset.X + squareSize/2.0, Y: y + g.offset.Y + squareSize/2.0}
				a      = geom.Vector2{X: x + g.offset.X, Y: y + g.offset.Y}
				b      = geom.Vector2{X: x + g.offset.X + squareSize, Y: y + g.offset.Y}
				c      = geom.Vector2{X: x + g.offset.X + squareSize, Y: y + g.offset.Y + squareSize}
				d      = geom.Vector2{X: x + g.offset.X, Y: y + g.offset.Y + squareSize}
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

			if g.isInsidePolygonWithHoles(center) {
				isCenter = true
			}

			// The square’s corners
			square := Square{
				A:        a,
				B:        b,
				C:        c,
				D:        d,
				Center:   center,
				isA:      isA,
				isB:      isB,
				isC:      isC,
				isD:      isD,
				isCenter: isCenter,
			}

			squares = append(squares, square)
			if square.isInside() {
				visSquares = append(visSquares, square)
			}
		}
	}

	return squares, visSquares
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

// isInsidePolygonWithHoles checks if p is inside the outer polygon but not inside any holes
func (g *Grid) isInsidePolygonWithHoles(point geom.Vector2) bool {
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

// isLineSegmentInsidePolygonOrHoles function to check if a line segment is inside a polygon or holes
func (g *Grid) isLineSegmentInsidePolygonOrHoles(lineStart, lineEnd geom.Vector2) bool {
	// Check if the line segment intersects any of the polygon edges
	if isLineSegmentInsidePolygon(g.polygon, lineStart, lineEnd) {
		return false
	}

	for _, hole := range g.holes {
		for i := 0; i < len(hole)-1; i++ {
			if doLinesIntersect(lineStart, lineEnd, hole[i], hole[i+1]) {
				return false // Intersects a hole
			}
		}
	}

	return true // Line segment is completely inside the polygon
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

// Function to check if two lines intersect
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

func heuristicEvaluation(a, b geom.Vector2) float32 {
	x := a.X - b.X
	y := a.Y - b.Y

	return float32(math.Sqrt(float64(x*x + y*y)))
}
