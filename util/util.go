package util

import (
	"log"
	"math/rand"

	"golang.org/x/exp/constraints"
)

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
