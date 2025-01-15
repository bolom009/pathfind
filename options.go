package pathfind

import (
	"github.com/bolom009/pathfind/graphs"
	"github.com/bolom009/pathfind/obstacles"
)

// PathOption represent setter of options to navigation
type PathOption func(*graphs.NavOpts)

// WithObstacles add obstacles to calc navigation
func WithObstacles(obstacles []obstacles.Obstacle) PathOption {
	return func(o *graphs.NavOpts) {
		o.Obstacles = obstacles
	}
}

// WithAgentRadius add agent radius to calc navigation
func WithAgentRadius(agentRadius float32) PathOption {
	return func(o *graphs.NavOpts) {
		o.AgentRadius = agentRadius
	}
}
