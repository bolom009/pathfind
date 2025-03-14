package recast

import (
	"github.com/bolom009/geom"
	"math"
	"sort"
)

// KDTree represents a KD-Tree node
type KDTree struct {
	Point geom.Vector2
	Left  *KDTree
	Right *KDTree
	Axis  int // 0 for X axis, 1 for Y axis
}

// BuildKDTree builds a KD-Tree from a slice of points
func BuildKDTree(points []geom.Vector2, depth int) *KDTree {
	if len(points) == 0 {
		return nil
	}

	// Alternate between X (0) and Y (1) axis
	axis := depth % 2
	sort.Slice(points, func(i, j int) bool {
		if axis == 0 {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})

	mid := len(points) / 2
	node := &KDTree{
		Point: points[mid],
		Axis:  axis,
	}

	node.Left = BuildKDTree(points[:mid], depth+1)
	node.Right = BuildKDTree(points[mid+1:], depth+1)

	return node
}

func (kd *KDTree) search(target geom.Vector2) (geom.Vector2, float32) {
	var (
		searchPoint geom.Vector2
		bestDest    = float32(math.Inf(1))
	)

	kd.nearestNeighborSearch(target, 0, &searchPoint, &bestDest)

	return searchPoint, bestDest
}

// nearestNeighborSearch searches for the nearest point to the target point in the KD-Tree
func (kd *KDTree) nearestNeighborSearch(target geom.Vector2, depth int, best *geom.Vector2, bestDist *float32) {
	if kd == nil {
		return
	}

	// Calculate current distance
	dist := geom.Distance(kd.Point, target)

	// Update best point and best distance
	if dist < *bestDist {
		*best = kd.Point
		*bestDist = dist
	}

	var (
		dp     = depth % 2
		dpZero = dp == 0
	)

	// Determine which side of the split to search first
	var nextBranch, otherBranch *KDTree
	if (dpZero && target.X < kd.Point.X) || (!dpZero && target.Y < kd.Point.Y) {
		nextBranch = kd.Left
		otherBranch = kd.Right
	} else {
		nextBranch = kd.Right
		otherBranch = kd.Left
	}

	// Search the next branch
	nextBranch.nearestNeighborSearch(target, depth+1, best, bestDist)

	// Check if we need to search the other branch
	if (dpZero && float32(math.Abs(float64(target.X-kd.Point.X))) < *bestDist) ||
		(!dpZero && float32(math.Abs(float64(target.Y-kd.Point.Y))) < *bestDist) {
		otherBranch.nearestNeighborSearch(target, depth+1, best, bestDist)
	}
}
