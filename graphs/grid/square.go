package grid

import (
	"github.com/bolom009/pathfind/vec"
)

// Square represent square with four points, also include info about each point if it's inside polygon
type Square struct {
	A, B, C, D, Center vec.Vector2
	// represent each point inside polygon
	isA, isB, isC, isD, isCenter bool
}

// Edges return list of all square edges
func (s *Square) Edges() []Edge {
	return []Edge{
		{A: vec.Vector2{X: s.A.X, Y: s.A.Y}, B: vec.Vector2{X: s.B.X, Y: s.B.Y}},
		{A: vec.Vector2{X: s.B.X, Y: s.B.Y}, B: vec.Vector2{X: s.C.X, Y: s.C.Y}},
		{A: vec.Vector2{X: s.C.X, Y: s.C.Y}, B: vec.Vector2{X: s.D.X, Y: s.D.Y}},
		{A: vec.Vector2{X: s.D.X, Y: s.D.Y}, B: vec.Vector2{X: s.A.X, Y: s.A.Y}},
	}
}

// Edge represent square edge
type Edge struct {
	A, B vec.Vector2
}

// isPointInsideSquare checks if point inside square
func (s *Square) isPointInsideSquare(point vec.Vector2) bool {
	// Assuming the square is axis-aligned, we can check the coordinates
	minX := s.A.X
	maxX := s.A.X
	minY := s.A.Y
	maxY := s.A.Y

	if s.B.X < minX {
		minX = s.B.X
	}
	if s.C.X < minX {
		minX = s.C.X
	}
	if s.D.X < minX {
		minX = s.D.X
	}

	if s.B.X > maxX {
		maxX = s.B.X
	}
	if s.C.X > maxX {
		maxX = s.C.X
	}
	if s.D.X > maxX {
		maxX = s.D.X
	}

	if s.B.Y < minY {
		minY = s.B.Y
	}
	if s.C.Y < minY {
		minY = s.C.Y
	}
	if s.D.Y < minY {
		minY = s.D.Y
	}

	if s.B.Y > maxY {
		maxY = s.B.Y
	}
	if s.C.Y > maxY {
		maxY = s.C.Y
	}
	if s.D.Y > maxY {
		maxY = s.D.Y
	}

	return point.X >= minX && point.X <= maxX && point.Y >= minY && point.Y <= maxY
}
