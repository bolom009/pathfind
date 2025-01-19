package modifiers

import (
	"math"

	"github.com/bolom009/geom"
)

// Simple smoothing a path by either moving the points closer together
func Simple(path []geom.Vector2, uniformLength bool, maxSegmentLength float32, subdivisions int, strength float32, iterations int) []geom.Vector2 {
	if len(path) < 2 {
		return path
	}

	subdivided := make([]geom.Vector2, 0, len(path))
	if uniformLength {
		if maxSegmentLength < 0.005 {
			maxSegmentLength = 0.005
		}

		var pathLength float32
		for i := 0; i < len(path)-1; i++ {
			pathLength += geom.Distance(path[i], path[i+1])
		}

		estimatedNumberOfSegments := int(math.Floor(float64(pathLength / maxSegmentLength)))
		subdivided = make([]geom.Vector2, 0, estimatedNumberOfSegments+2)

		var distanceAlong float32
		for i := 0; i < len(path)-1; i++ {
			var (
				start  = path[i]
				end    = path[i+1]
				length = geom.Distance(start, end)
			)
			for distanceAlong < length {
				subdivided = append(subdivided, start.Lerp(end, distanceAlong/length))
				distanceAlong += maxSegmentLength
			}

			distanceAlong -= length
		}

		subdivided = append(subdivided, path[len(path)-1])
	} else {
		if subdivisions < 0 {
			subdivisions = 0
		}

		if subdivisions > 10 {
			subdivisions = 10
		}

		subSegments := 1 << subdivisions
		subdivided = subdivide(path, subSegments)
	}

	if strength > 0 {
		for it := 0; it < iterations; it++ {
			prev := subdivided[0]
			for i := 1; i < len(subdivided)-1; i++ {
				tmp := subdivided[i]
				subdivided[i] = tmp.Lerp(geom.Vector2{
					X: (prev.X + subdivided[i+1].X) / 2,
					Y: (prev.Y + subdivided[i+1].Y) / 2,
				}, strength)

				prev = tmp
			}
		}
	}

	return subdivided
}

func subdivide(path []geom.Vector2, subSegments int) []geom.Vector2 {
	result := make([]geom.Vector2, 0)
	for i := 0; i < len(path)-1; i++ {
		for j := 0; j < subSegments; j++ {
			point := path[i].Lerp(path[i+1], float32(j/subSegments))
			result = append(result, point)
		}
	}

	result = append(result, path[len(path)-1])
	return result
}
