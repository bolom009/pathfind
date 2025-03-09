package recast

import (
	"github.com/bolom009/geom"
	"testing"
)

func BenchmarkTriangle_pointInsideTriangle(b *testing.B) {
	p1 := geom.Vector2{X: 14.142136, Y: 14.142136}
	p2 := geom.Vector2{X: 105.857864, Y: 14.142136}
	p3 := geom.Vector2{X: 105.857864, Y: 354.14215}
	pp := geom.Vector2{X: 25, Y: 25}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 18; j++ {
			_ = pointInsideTriangle(p1, p2, p3, pp)
		}
	}
}
