package recast

import (
	"github.com/bolom009/geom"
	"github.com/fzipp/astar"
)

type option func(r *Recast)

func WithCostFunc(costFunc astar.CostFunc[geom.Vector2]) option {
	return func(g *Recast) {
		g.costFunc = costFunc
	}
}
