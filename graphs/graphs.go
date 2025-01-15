package graphs

import (
	"context"
	"iter"
	"slices"

	"github.com/bolom009/pathfind/obstacles"
)

// NavGraph is represented an interface for graph types
type NavGraph[Node comparable] interface {
	Generate(ctx context.Context) error
	AggregationGraph(Node, Node, []obstacles.Obstacle) Graph[Node]
	GetVisibility() Graph[Node]
	ContainsPoint(Node) bool
	Cost(Node, Node) float64
}

// Graph is represented by an adjacency list.
type Graph[Node comparable] map[Node][]Node

// Link creates a directed edge from node a to node b.
func (g Graph[Node]) Link(a, b Node) Graph[Node] {
	g[a] = append(g[a], b)
	return g
}

func (g Graph[Node]) DeleteNode(node Node) Graph[Node] {
	if _, ok := g[node]; ok {
		delete(g, node)
	}

	return g
}

func (g Graph[Node]) DeleteNeighbour(node, neighbour Node) Graph[Node] {
	if neighbours, ok := g[node]; ok {
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
	cGraph := make(Graph[Node])
	for k, v := range g {
		cv := make([]Node, len(v))
		copy(cv, v)

		cGraph[k] = cv
	}
	return cGraph
}

// Neighbours returns the neighbour nodes of node n in the graph.
// This method makes graph[Node] implement the astar.Graph[Node] interface.
func (g Graph[Node]) Neighbours(n Node) iter.Seq[Node] {
	return slices.Values(g[n])
}
