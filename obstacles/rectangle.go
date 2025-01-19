package obstacles

import (
	"github.com/bolom009/geom"
	"math"
)

// Rectangle is represented obstacle in rectangle shape
type Rectangle struct {
	width   float32
	height  float32
	polygon []geom.Vector2
	center  geom.Vector2
}

func (o *Rectangle) GetCenter() geom.Vector2 {
	return o.center
}

func (o *Rectangle) GetPolygon() []geom.Vector2 {
	return o.polygon
}

func (o *Rectangle) Move(pos geom.Vector2) {
	o.center = o.center.Add(pos)
	for i := 0; i < len(o.polygon); i++ {
		o.polygon[i] = o.polygon[i].Add(pos)
	}
}

func (o *Rectangle) IsPointAround(point geom.Vector2, edgeLen float32) bool {
	distance := edgeLen
	halfWidth := o.width / 2
	halfHeight := o.height / 2

	// Calculate the boundaries of the rectangle
	left := o.center.X - halfWidth
	right := o.center.X + halfWidth
	bottom := o.center.Y - halfHeight
	top := o.center.Y + halfHeight

	// Check if the point is within the rectangle
	if point.X >= left && point.X <= right && point.Y >= bottom && point.Y <= top {
		return true // The point is inside the rectangle
	}

	// Calculate the closest distance from the point to the rectangle
	closestX := math.Max(float64(left), math.Min(float64(point.X), float64(right)))
	closestY := math.Max(float64(bottom), math.Min(float64(point.Y), float64(top)))

	// Calculate the Euclidean distance from the point to the closest point on the rectangle
	distanceToClosestPoint := math.Sqrt(math.Pow(float64(point.X)-closestX, 2) + math.Pow(float64(point.Y)-closestY, 2))

	// Check if the distance is within the specified range
	return distanceToClosestPoint <= float64(distance)
}

func GenerateRectangle(center geom.Vector2, width, height float32) *Rectangle {
	var (
		halfWidth  = width / 2
		halfHeight = height / 2
	)

	return &Rectangle{
		width:  width,
		height: height,
		polygon: []geom.Vector2{
			{X: center.X - halfWidth, Y: center.Y - halfHeight},
			{X: center.X + halfWidth, Y: center.Y - halfHeight},
			{X: center.X + halfWidth, Y: center.Y + halfHeight},
			{X: center.X - halfWidth, Y: center.Y + halfHeight},
		},
		center: center,
	}
}
