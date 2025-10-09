package recast

import (
	"math"

	goclipper2 "github.com/bolom009/go-clipper2"
	"github.com/bolom009/pathfind/mesh"
)

func (r *Recast) getPolyOffsetsWithUnion(polygon *mesh.Polygon) goclipper2.PathsD {
	subject := make(goclipper2.PathsD, 0)
	newPaths := inflatePathsD(goclipper2.PathsD{toPathD(polygon.Points())}, float64(-polygon.Offset()), goclipper2.Miter, goclipper2.Polygon)
	subject = append(subject, newPaths...)

	for _, obstacle := range polygon.Obstacles() {
		rPoints := goclipper2.ReversePath(obstacle.Points())
		newPaths := inflatePathsD(goclipper2.PathsD{toPathD(rPoints)}, float64(obstacle.Offset()), goclipper2.Miter, goclipper2.Polygon)

		subject = append(subject, newPaths...)
	}

	for _, innerHole := range polygon.InnerHoles() {
		newPaths := inflatePathsD(goclipper2.PathsD{toPathD(innerHole.Points())}, float64(innerHole.Offset()), goclipper2.Miter, goclipper2.Polygon)

		subject = append(subject, newPaths...)
	}

	offsetPath := unionPathsD(subject, goclipper2.Positive)
	return offsetPath
}

// inflatePathsD override standard method without decimal calc
func inflatePathsD(paths goclipper2.PathsD, delta float64, joinType goclipper2.JoinType, endType goclipper2.EndType) goclipper2.PathsD {
	miterLimit := 2.0
	arcTolerance := 0.0
	precision := 2
	scale := math.Pow(10, float64(precision))
	tmp := scalePathsDToPaths64(paths, scale)

	co := goclipper2.NewClipperOffset(miterLimit, scale*arcTolerance, false, false)
	co.AddPaths(tmp, joinType, endType)
	co.Execute64(delta*scale, &tmp)

	return scalePaths64ToPathsD(tmp, 1.0/scale)
}

func scalePathsDToPaths64(paths goclipper2.PathsD, scale float64) goclipper2.Paths64 {
	result := make(goclipper2.Paths64, len(paths))
	for i, path := range paths {
		result[i] = scalePathDToPath64(path, scale)
	}

	return result
}

func scalePathDToPath64(path goclipper2.PathD, scale float64) goclipper2.Path64 {
	result := make(goclipper2.Path64, len(path))
	for i, pt := range path {
		result[i] = goclipper2.Point64{X: int64(pt.X * scale), Y: int64(pt.Y * scale)}
	}

	return result
}

func scalePaths64ToPathsD(paths goclipper2.Paths64, scale float64) goclipper2.PathsD {
	result := make(goclipper2.PathsD, len(paths))
	for i, path := range paths {
		result[i] = scalePath64ToPathD(path, scale)
	}

	return result
}

func scalePath64ToPathD(path goclipper2.Path64, scale float64) goclipper2.PathD {
	result := make(goclipper2.PathD, len(path))
	for i, pt := range path {
		result[i] = goclipper2.PointD{X: float64(pt.X) * scale, Y: float64(pt.Y) * scale}
	}

	return result
}

func unionPathsD(subject goclipper2.PathsD, fillRule goclipper2.FillRule) goclipper2.PathsD {
	solution := make(goclipper2.PathsD, 0)
	c := goclipper2.NewClipperD(2)
	c.AddPathsWithScaleFunc(subject, goclipper2.Subject, false, scalePathsDToPaths64)

	c.ExecuteWithScaleFunc(goclipper2.Union, fillRule, &solution, nil, scalePath64ToPathD)
	return solution
}
