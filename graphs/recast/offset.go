package recast

import (
	goclipper2 "github.com/bolom009/go-clipper2"
	"github.com/bolom009/pathfind/mesh"
)

func getPolyOffsetsWithUnion(polygon *mesh.Polygon) goclipper2.PathsD {
	subject := make(goclipper2.PathsD, 0)
	newPaths := goclipper2.InflatePathsD(goclipper2.PathsD{toPathD(polygon.Points())}, float64(-polygon.Offset()), goclipper2.Miter, goclipper2.Polygon)
	subject = append(subject, newPaths...)

	for _, obstacle := range polygon.Obstacles() {
		rPoints := goclipper2.ReversePath(obstacle.Points())
		newPaths := goclipper2.InflatePathsD(goclipper2.PathsD{toPathD(rPoints)}, float64(obstacle.Offset()), goclipper2.Miter, goclipper2.Polygon)

		subject = append(subject, newPaths...)
	}

	for _, innerHole := range polygon.InnerHoles() {
		newPaths := goclipper2.InflatePathsD(goclipper2.PathsD{toPathD(innerHole.Points())}, float64(innerHole.Offset()), goclipper2.Miter, goclipper2.Polygon)

		subject = append(subject, newPaths...)
	}

	offsetPath := goclipper2.UnionPathsD(subject, goclipper2.Positive)
	return offsetPath
}
