package pathfind

import (
	"context"
	"fmt"

	"github.com/bolom009/astar"
	"github.com/bolom009/pathfind/graphs"
)

// Pathfinder represent struct to operate with pathfinding
// Graph of current pathfinder represent as list of squares
type Pathfinder[Node comparable] struct {
	graphs []graphs.NavGraph[Node]
}

// NewPathfinder constructor to create pathfinder struct
// polygons contains two parts: polygon and holes
func NewPathfinder[Node comparable](graphs []graphs.NavGraph[Node]) *Pathfinder[Node] {
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
// To search path could be added dynamic obstacles. All obstacles will cut current graph by their polygon.
// The function returns nil if no path exists
func (p *Pathfinder[Node]) Path(graphID int, start, dest Node, opts ...PathOption) []Node {
	g := p.graphs[graphID]
	if g == nil {
		return nil
	}

	navOpts := &graphs.NavOpts{}
	for _, opt := range opts {
		opt(navOpts)
	}

	vis := g.AggregationGraph(start, dest, nil)

	// check if start-dest graph then skip A* search path
	if len(vis) == 2 && vis[start][0] == dest && vis[dest][0] == start {
		return []Node{start, dest}
	}

	path := astar.FindPath[Node](vis, start, dest, g.Cost, g.Cost)

	return path
}

// Graph return generated graph visibility
func (p *Pathfinder[Node]) Graph(graphID int, opts ...PathOption) map[Node][]Node {
	g := p.graphs[graphID]
	if g == nil {
		return nil
	}

	navOpts := &graphs.NavOpts{}
	for _, opt := range opts {
		opt(navOpts)
	}

	return g.GetVisibility(navOpts)
}

func (p *Pathfinder[Node]) GraphsNum() int {
	return len(p.graphs)
}

func (p *Pathfinder[Node]) GetClosestPoint(graphID int, point Node) (Node, bool) {
	g := p.graphs[graphID]
	if g == nil {
		var zero Node
		return zero, false
	}

	return g.GetClosestPoint(point)
}

func (p *Pathfinder[Node]) IsRaycastHit(graphID int, start, dest Node) bool {
	g := p.graphs[graphID]
	if g == nil {
		return false
	}

	return g.IsRaycastHit(start, dest)
}

// GraphWithSearchPath return generated graph with path nodes
func (p *Pathfinder[Node]) GraphWithSearchPath(graphID int, start, dest Node, opts ...PathOption) graphs.Graph[Node] {
	g := p.graphs[graphID]
	if g == nil {
		return nil
	}

	navOpts := &graphs.NavOpts{}
	for _, opt := range opts {
		opt(navOpts)
	}

	return g.AggregationGraph(start, dest, navOpts)
}
