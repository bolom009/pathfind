package mesh

import (
	"github.com/bolom009/clipper"
	"github.com/bolom009/geom"
)

type Polygon struct {
	points []geom.Vector2
	offset float32
}

func NewPolygon(points []geom.Vector2, offset float32) *Polygon {
	return &Polygon{
		points: points,
		offset: offset,
	}
}

func (p *Polygon) Points() []geom.Vector2 {
	return p.points
}

func (p *Polygon) Offset() float32 {
	return p.offset
}

func (p *Polygon) Clipper() []geom.Vector2 {
	return clipper.OffsetPolygon(p.points, p.offset)
}
