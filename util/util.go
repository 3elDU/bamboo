package util

import (
	"log"
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

// Limits float precision to `points` digits after comma
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
		log.Panicln("array with zero length")
	}
	return objects[rand.Intn(len(objects))]
}
