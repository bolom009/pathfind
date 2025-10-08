package mesh

import (
	"github.com/bolom009/geom"
)

type HoleType uint8

const (
	InnerHole HoleType = iota + 1
	Obstacle
)

type Hole struct {
	points   []geom.Vector2
	hType    HoleType
	offset   float32
	viewable bool
}

func NewInnerHole(points []geom.Vector2, offset float32) *Hole {
	return &Hole{
		points:   points,
		offset:   offset,
		viewable: false,
		hType:    Obstacle,
	}
}

func NewObstacle(points []geom.Vector2, offset float32, viewable bool) *Hole {
	return &Hole{
		points:   points,
		offset:   offset,
		viewable: viewable,
		hType:    Obstacle,
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

func (h *Hole) Type() HoleType {
	return h.hType
}
