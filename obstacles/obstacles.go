package obstacles

import "github.com/bolom009/geom"

// Obstacle is represented an interface for shape of obstacle
type Obstacle interface {
	GetCenter() geom.Vector2
	GetPolygon() []geom.Vector2
	Move(vector2 geom.Vector2)
	IsPointAround(point geom.Vector2, edgeLen float32) bool
}
