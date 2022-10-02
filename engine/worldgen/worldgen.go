package worldgen

import (
	"github.com/aquilax/go-perlin"
	"github.com/veandco/go-sdl2/sdl"
)

type World struct {
	generator   *perlin.Perlin
	seed        int64
	ScaleFactor float64
}

func New(seed int64, scaleFactor float64) *World {
	return &World{
		generator:   perlin.NewPerlin(2, 2, 3, seed),
		seed:        seed,
		ScaleFactor: scaleFactor,
	}
}

func (w *World) Block(x, y float64) sdl.Color {
	h := w.generator.Noise2D(x/w.ScaleFactor, y/w.ScaleFactor) + 1

	switch {
	case h <= 1:
		return sdl.Color{R: 0, G: 0, B: 255, A: 255}
	case h <= 1.1:
		return sdl.Color{R: 255, G: 255, B: 0, A: 255}
	case h <= 1.6:
		return sdl.Color{R: 0, G: 255, B: 0, A: 255}
	case h <= 1.8:
		return sdl.Color{R: 128, G: 128, B: 128, A: 255}
	default:
		return sdl.Color{R: 255, G: 255, B: 255, A: 255}
	}
}
