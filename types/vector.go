package types

import "math"

// A generic vector interface that contains two float64 values
type IVector interface {
	Values() (float64, float64)
	DistanceTo(other IVector) float64
}

type Vec2i struct {
	X, Y int
}

func (vec Vec2i) Values() (float64, float64) {
	return float64(vec.X), float64(vec.Y)
}

func (vec Vec2i) DistanceTo(other IVector) float64 {
	x2, y2 := other.Values()
	return math.Abs(float64(vec.X)-x2) + math.Abs(float64(vec.Y)-y2)
}

type Vec2u struct {
	X, Y uint64
}

func (vec Vec2u) Values() (float64, float64) {
	return float64(vec.X), float64(vec.Y)
}

func (vec Vec2u) DistanceTo(other IVector) float64 {
	x2, y2 := other.Values()
	return math.Abs(float64(vec.X)-x2) + math.Abs(float64(vec.Y)-y2)
}

type Vec2f struct {
	X, Y float64
}

func (vec Vec2f) Values() (float64, float64) {
	return vec.X, vec.Y
}
func (vec Vec2f) DistanceTo(other IVector) float64 {
	x2, y2 := other.Values()
	return math.Abs(vec.X-x2) + math.Abs(vec.Y-y2)
}
