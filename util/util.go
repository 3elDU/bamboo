package util

import (
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

type Coords2i struct {
	X, Y int64
}

type Coords2u struct {
	X, Y uint64
}

type Coords2f struct {
	X, Y float64
}

// Limits float precision to `points` after comma
func LimitFloatPrecision(val float64, points int) float64 {
	return math.Round(val*float64(points)*10) / (float64(points) * 10)
}

func Clamp[T constraints.Integer | constraints.Float](val, min, max T) T {
	if val < min {
		return min
	} else if val > max {
		return max
	} else {
		return val
	}
}

func RandomChoice[T any](objects []T) T {
	if len(objects) == 0 {
		panic("objects with zero length")
	}
	return objects[rand.Intn(len(objects))]
}

// Convert block coordinates to chunk coordinates
func CoordsBlockToChunk(blockCoords Coords2i) Coords2i {
	blockCoords.X += 8
	blockCoords.Y += 8
	return Coords2i{X: blockCoords.X / 16, Y: blockCoords.Y / 16}
}
