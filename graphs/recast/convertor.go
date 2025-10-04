package recast

import (
	"github.com/bolom009/geom"
	goclipper2 "github.com/bolom009/go-clipper2"
	"github.com/govalues/decimal"
)

func toPoint2(paths goclipper2.PathD) []geom.Vector2 {
	res := make([]geom.Vector2, len(paths))
	for i, path := range paths {
		res[i] = geom.Vector2{
			X: float32(path.X),
			Y: float32(path.Y),
		}
	}

	return res
}

func toPathD(paths []geom.Vector2) goclipper2.PathD {
	res := make(goclipper2.PathD, len(paths))
	for i, path := range paths {
		xV := decimal.Decimal{}
		_ = xV.Scan(path.X)

		yV := decimal.Decimal{}
		_ = yV.Scan(path.Y)

		x, _ := xV.Float64()
		y, _ := yV.Float64()

		res[i] = goclipper2.PointD{
			X: x,
			Y: y,
		}
	}

	return res
}
