package graphs

import (
	"context"
	"github.com/bolom009/pathfind/obstacles"
	"maps"
)

// NavOpts contains optional objects for aggregation graph
// Obstacles parameter will exclude nodes & edges of generated graph by obstacles polygon
// AgentRadius parameter will take into account the path search by agent radius
type NavOpts struct {
	Obstacles   []obstacles.Obstacle
	AgentRadius float32
}

// NavGraph is represented an interface for graph types
type NavGraph[Node comparable] interface {
	Generate(ctx context.Context) error
	AggregationGraph(Node, Node, *NavOpts) Graph[Node]
	GetVisibility(opts *NavOpts) Graph[Node]
	ContainsPoint(Node) bool
	Cost(Node, Node) float32
}

// Graph is represented by an adjacency list.
type Graph[Node comparable] map[Node][]Node

// Link creates a directed edge from node a to node b.
func (g Graph[Node]) Link(a, b Node) Graph[Node] {
	g[a] = append(g[a], b)
	return g
}

// LinkBoth creates a both directed edge from node a to node b and back
func (g Graph[Node]) LinkBoth(a, b Node, l ...int32) Graph[Node] {
	if len(l) > 0 {
		if _, ok := g[a]; !ok {
			g[a] = make([]Node, 0, l[0])
			g[a] = append(g[a], b)
		} else {
			g[a] = append(g[a], b)
		}

		if _, ok := g[b]; !ok {
			g[b] = make([]Node, 0, l[0])
			g[b] = append(g[b], a)
		} else {
			g[b] = append(g[b], a)
		}
	} else {
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}

	return g
}

// DeleteNode delete graph vertex if exist
func (g Graph[Node]) DeleteNode(node Node) Graph[Node] {
	if _, ok := g[node]; ok {
		delete(g, node)
	}

	return g
}

// DeleteNeighbour delete vertex edge if exist
func (g Graph[Node]) DeleteNeighbour(node, neighbour Node) Graph[Node] {
	if neighbours, ok := g[node]; ok {
		if len(neighbours) == 0 {
			return g
		}

		newSlice := make([]Node, 0, len(neighbours)-1)
		for _, n := range neighbours {
			if neighbour != n {
				newSlice = append(newSlice, n)
			}
		}

		g[node] = newSlice
	}

	return g
}

// Copy return copy of graph
func (g Graph[Node]) Copy() Graph[Node] {
	return maps.Clone(g)
}

// Neighbours returns the neighbour nodes of node n in the graph.
// This method makes graph[Node] implement the astar.Graph[Node] interface.
func (g Graph[Node]) Neighbours(n Node) []Node {
	return g[n]
}
