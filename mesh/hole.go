package mesh

import (
	"github.com/bolom009/geom"
)

type Hole struct {
	points   []geom.Vector2
	offset   float32
	viewable bool
}

func NewInnerHole(points []geom.Vector2, offset float32) *Hole {
	return &Hole{
		points:   points,
		offset:   offset,
		viewable: false,
	}
}

func NewObstacle(points []geom.Vector2, offset float32, viewable bool) *Hole {
	return &Hole{
		points:   points,
		offset:   offset,
		viewable: viewable,
	}
}

func (h *Hole) Points() []geom.Vector2 {
	return h.points
}

func (h *Hole) Offset() float32 {
	return h.offset
}

func (h *Hole) Viewable() bool {
	return h.viewable
}
