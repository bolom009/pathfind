package pathfind

import "math"

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector2) Sub(q Vector2) Vector2 {
	return Vector2{v.X - q.X, v.Y - q.Y}
}

func (v Vector2) Scale(s float64) Vector2 {
	return Vector2{X: v.X * s, Y: v.Y * s}
}

func (v Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v Vector2) Div(s float64) Vector2 {
	return Vector2{X: v.X / s, Y: v.Y / s}
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) Normalize() Vector2 {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2{0, 0}
	}
	return Vector2{v.X / mag, v.Y / mag}
}

func Distance(a, b Vector2) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}

func Lerp(a, b Vector2, t float64) Vector2 {
	return Vector2{X: a.X + (b.X-a.X)*t, Y: a.Y + (b.Y-a.Y)*t}
}
