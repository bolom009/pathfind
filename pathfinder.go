package pathfind

import (
	"context"
	"fmt"
	"github.com/bolom009/pathfind/graphs"
	"github.com/fzipp/astar"
)

// Pathfinder represent struct to operate with pathfinding
// Graph of current pathfinder represent as list of squares
type Pathfinder[Node comparable] struct {
	graphs []graphs.IGraph[Node]
}

// NewPathfinder constructor to create pathfinder struct
// polygons contains two parts: polygon and holes
func NewPathfinder[Node comparable](graphs []graphs.IGraph[Node]) *Pathfinder[Node] {
	p := &Pathfinder[Node]{
		graphs: graphs,
	}

	return p
}

// Initialize split polygon with holes to square and prepare visibility graph
func (p *Pathfinder[Node]) Initialize(ctx context.Context) error {
	for _, g := range p.graphs {
		if err := g.Generate(ctx); err != nil {
			return fmt.Errorf("graph generate: %w", err)
		}
	}

	return nil
}

// Path finds the shortest path from start to dest
// If dest is outside the polygon set it will be clamped to the nearest polygon point.
// The function returns nil if no path exists
func (p *Pathfinder[Node]) Path(graphID int, start, dest Node) []Node {
	g := p.graphs[graphID]
	if g == nil {
		return nil
	}

	vis := g.AggregationGraph(start, dest)

	return astar.FindPath[Node](vis, start, dest, g.Cost, g.Cost)
}

// Graph return generated graph based on square list
func (p *Pathfinder[Node]) Graph(graphID int) map[Node][]Node {
	g := p.graphs[graphID]
	if g == nil {
		return nil
	}

	return g.GetVisibility()
}

func (p *Pathfinder[Node]) GraphsNum() int {
	return len(p.graphs)
}

// GraphWithSearchPath return generated graph with path nodes
func (p *Pathfinder[Node]) GraphWithSearchPath(graphID int, start, dest Node) map[Node][]Node {
	g := p.graphs[graphID]
	if g == nil {
		return nil
	}

	return g.AggregationGraph(start, dest)
}
