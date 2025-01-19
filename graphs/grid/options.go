package grid

import (
	"github.com/bolom009/geom"
	"github.com/fzipp/astar"
)

type option func(g *Grid)

func WithCostFunc(costFunc astar.CostFunc[geom.Vector2]) option {
	return func(g *Grid) {
		g.costFunc = costFunc
	}
}

func WithOffset(offset geom.Vector2) option {
	return func(g *Grid) {
		g.offset = offset
	}
}
