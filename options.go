package pathfind

import "github.com/fzipp/astar"

type Option func(p *Pathfinder)

func WithCostFunc(constFunc astar.CostFunc[Vector2]) Option {
	return func(p *Pathfinder) {
		p.constFunc = constFunc
	}
}
