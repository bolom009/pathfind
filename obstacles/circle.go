package obstacles

import (
	"github.com/bolom009/pathfind/vec"
	"math"
)

// Circle is represent obstacle in circle shape
type Circle struct {
	radius   float32
	segments int
	polygon  []vec.Vector2
	center   vec.Vector2
}

func (o *Circle) GetCenter() vec.Vector2 {
	return o.center
}

func (o *Circle) GetPolygon() []vec.Vector2 {
	return o.polygon
}

func (o *Circle) Move(pos vec.Vector2) {
	o.center = o.center.Add(pos)
	for i := 0; i < len(o.polygon); i++ {
		o.polygon[i] = o.polygon[i].Add(pos)
	}
}

func (o *Circle) IsPointAround(point vec.Vector2, edgeLen float32) bool {
	rangeSquared := edgeLen * edgeLen
	dx := point.X - o.center.X
	dy := point.Y - o.center.Y
	distanceSquared := dx*dx + dy*dy
	radiusSquared := o.radius * o.radius

	// Check if the point is inside the circle
	if distanceSquared <= radiusSquared {
		return true // The point is inside the circle
	}

	// Check if the point is close to the circle
	return distanceSquared <= (radiusSquared + rangeSquared)
}

func GenerateCircle(center vec.Vector2, radius float32, segments int) *Circle {
	if segments < 3 {
		panic("a polygon needs at least 3 segments (like a triangle)")
	}

	points := make([]vec.Vector2, segments)
	angleIncrement := 2 * math.Pi / float64(segments) // Full circle divided by number of segments

	for i := 0; i < segments; i++ {
		angle := angleIncrement * float64(i)
		x := center.X + radius*float32(math.Cos(angle))
		y := center.Y + radius*float32(math.Sin(angle))
		points[i] = vec.Vector2{X: x, Y: y}
	}

	return &Circle{
		polygon:  points,
		center:   center,
		radius:   radius,
		segments: segments,
	}
}
