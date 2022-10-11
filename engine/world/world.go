package world

import (
	"fmt"
	"image/color"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/util"
	"github.com/aquilax/go-perlin"
)

type World struct {
	generator *perlin.Perlin
	seed      int64

	// keys there are chunk coordinates.
	// so, actual chunk coordinates are x*16 and y*16
	data map[util.Coords2i]*chunk
}

func NewWorld(seed int64) *World {
	return &World{
		generator: perlin.NewPerlin(2, 2, 16, seed),
		seed:      seed,
		data:      make(map[util.Coords2i]*chunk),
	}
}

func (w *World) gen(x, y float64) color.RGBA {
	h := w.generator.Noise2D(x/config.PERLIN_NOISE_SCALE_FACTOR, y/config.PERLIN_NOISE_SCALE_FACTOR) + 1

	switch {
	case h <= 1:
		return colors.DarkBlue
	case h <= 1.1:
		return colors.Yellow
	case h <= 1.6:
		return colors.Green
	case h <= 1.8:
		return colors.DarkGray1
	default:
		return colors.Gray
	}
}

func (w *World) GenerateChunk(cx, cy int64) {
	fmt.Printf("Generating chunk %v, %v\n", cx, cy)

	chunk := NewChunk(cx, cy)
	w.data[util.Coords2i{X: cx, Y: cy}] = chunk

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			clr := w.gen(float64(int(cx)*16+x), float64(int(cy)*16+y))
			chunk.SetBlock(x, y, NewColoredBlock(clr))
		}
	}
}

// x and y are player coordinates
func (w *World) Update(playerX, playerY float64) {

}

// Returns a chunk at given coordinates
// Note that x and y are chunk coordinates, not block coordinates
func (w *World) At(x, y int64) *chunk {
	_, exists := w.data[util.Coords2i{X: x, Y: y}]

	// generate chunk, if it doesn't exist yet
	if !exists {
		w.GenerateChunk(x, y)
	}

	return w.data[util.Coords2i{X: x, Y: y}]
}
