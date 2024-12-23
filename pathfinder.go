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
	graph graphs.IGraph[Node]
}

// NewPathfinder constructor to create pathfinder struct
// polygons contains two parts: polygon and holes
func NewPathfinder[Node comparable](graph graphs.IGraph[Node]) *Pathfinder[Node] {
	p := &Pathfinder[Node]{
		graph: graph,
	}

	return p
}

// Initialize split polygon with holes to square and prepare visibility graph
func (p *Pathfinder[Node]) Initialize(ctx context.Context) error {
	if err := p.graph.Generate(ctx); err != nil {
		return fmt.Errorf("graph generate: %w", err)
	}
	return nil
}

// Path finds the shortest path from start to dest
// If dest is outside the polygon set it will be clamped to the nearest polygon point.
// The function returns nil if no path exists
func (p *Pathfinder[Node]) Path(start, dest Node) []Node {
	vis := p.graph.AggregationGraph(start, dest)

	return astar.FindPath[Node](vis, start, dest, p.graph.Cost, p.graph.Cost)
}

// Graph return generated graph based on square list
func (p *Pathfinder[Node]) Graph() map[Node][]Node {
	return p.graph.GetVisibility()
}

// GraphWithSearchPath return generated graph with path nodes
func (p *Pathfinder[Node]) GraphWithSearchPath(start, dest Node) map[Node][]Node {
	return p.graph.AggregationGraph(start, dest)
}
