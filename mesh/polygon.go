package mesh

import (
	"github.com/bolom009/geom"
)

type Polygon struct {
	points []geom.Vector2
	// innerHoles represent holes inside polygon or holes created by closed circuit polygons
	innerHoles []*Hole
	// obstacles represent general obstacles like: wall, box, lake, building, arch etc.
	obstacles []*Hole
	// contains list of all polygon innerHoles and obstacles
	holes []*Hole
	// use for clipper2 offset process
	offset float32
}

func NewPolygon(points []geom.Vector2, innerHoles []*Hole, obstacles []*Hole, offset float32) *Polygon {
	return &Polygon{
		points:     points,
		offset:     offset,
		innerHoles: innerHoles,
		obstacles:  obstacles,
		holes:      append(obstacles, innerHoles...),
	}
}

func (p *Polygon) Points() []geom.Vector2 {
	return p.points
}

func (p *Polygon) Holes() []*Hole {
	return p.holes
}

func (p *Polygon) InnerHoles() []*Hole {
	return p.innerHoles
}

func (p *Polygon) Obstacles() []*Hole {
	return p.obstacles
}

func (p *Polygon) Offset() float32 {
	return p.offset
}
