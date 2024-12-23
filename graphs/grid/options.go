package grid

import (
	"github.com/bolom009/pathfind/vec"
	"github.com/fzipp/astar"
)

type option func(p *Grid)

func WithCostFunc(costFunc astar.CostFunc[vec.Vector2]) option {
	return func(p *Grid) {
		p.costFunc = costFunc
	}
}
