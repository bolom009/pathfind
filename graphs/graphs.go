package graphs

import (
	"context"
	"iter"
	"slices"
)

// IGraph is represented an interface for graph types
type IGraph[Node comparable] interface {
	Generate(ctx context.Context) error
	AggregationGraph(Node, Node) Graph[Node]
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

// Copy return copy of graph
func (g Graph[Node]) Copy() Graph[Node] {
	cGraph := make(Graph[Node])
	for k, v := range g {
		cGraph[k] = v
	}
	return cGraph
}

// Neighbours returns the neighbour nodes of node n in the graph.
// This method makes graph[Node] implement the astar.Graph[Node] interface.
func (g Graph[Node]) Neighbours(n Node) iter.Seq[Node] {
	return slices.Values(g[n])
}
