package obstacles

import "github.com/bolom009/pathfind/vec"

// Obstacle is represented an interface for shape of obstacle
type Obstacle interface {
	GetCenter() vec.Vector2
	GetPolygon() []vec.Vector2
	Move(vector2 vec.Vector2)
	IsPointAround(point vec.Vector2, edgeLen float32) bool
}
