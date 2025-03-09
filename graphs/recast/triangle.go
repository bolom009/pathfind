package recast

import (
	"github.com/bolom009/delaunay"
	"github.com/bolom009/geom"
)

type Triangle []geom.Vector2

func convertTriangles(dTriangles []*delaunay.Triangle) []Triangle {
	triangles := make([]Triangle, len(dTriangles))
	for i, triangle := range dTriangles {
		points := triangle.GetPoints()
		triangles[i] = Triangle{
			geom.Vector2{X: points[0].X, Y: points[0].Y},
			geom.Vector2{X: points[1].X, Y: points[1].Y},
			geom.Vector2{X: points[2].X, Y: points[2].Y},
		}
	}

	return triangles
}

func pointInsideTriangle(p1, p2, p3, pp geom.Vector2) bool {
	if area(p1, p2, p3) >= 0 {
		return area(p1, p2, pp) >= 0 && area(p2, p3, pp) >= 0 && area(p3, p1, pp) >= 0
	}

	return area(p1, p2, pp) <= 0 && area(p2, p3, pp) <= 0 && area(p3, p1, pp) <= 0
}

func area(p1, p2, p3 geom.Vector2) float32 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p1.Y-p3.Y)*(p2.X-p3.X)
}
