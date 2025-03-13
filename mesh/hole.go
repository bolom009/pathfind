package mesh

import (
	"github.com/bolom009/clipper"
	"github.com/bolom009/geom"
)

type Hole struct {
	points []geom.Vector2
	offset float32
}

func NewHole(points []geom.Vector2, offset float32) *Hole {
	return &Hole{
		points: points,
		offset: offset,
	}
}

func (h *Hole) Points() []geom.Vector2 {
	return h.points
}

func (h *Hole) Offset() float32 {
	return h.offset
}

func (h *Hole) Clipper() []geom.Vector2 {
	return clipper.OffsetPolygon(h.points, h.offset)
}
