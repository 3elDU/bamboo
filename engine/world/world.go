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

	// keys there are Chunk coordinates.
	// so, actual Chunk coordinates are x*16 and y*16
	data map[util.Coords2i]*Chunk
}

func NewWorld(seed int64) *World {
	return &World{
		generator: perlin.NewPerlin(2, 2, 16, seed),
		seed:      seed,
		data:      make(map[util.Coords2i]*Chunk),
	}
}

func (world *World) gen(x, y float64) color.RGBA {
	h := world.generator.Noise2D(x/config.PerlinNoiseScaleFactor, y/config.PerlinNoiseScaleFactor) + 1

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

func (world *World) GenerateChunk(cx, cy int64) error {
	fmt.Printf("Generating Chunk %v, %v\n", cx, cy)

	chunk := NewChunk(cx, cy)
	world.data[util.Coords2i{X: cx, Y: cy}] = chunk

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			clr := world.gen(float64(int(cx)*16+x), float64(int(cy)*16+y))
			err := chunk.SetBlock(x, y, NewColoredBlock(clr))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Update - x and y are player coordinates
func (world *World) Update(_, _ float64) {

}

// At Returns a Chunk at given coordinates. Note that x and y are Chunk
// coordinates, not block coordinates
func (world *World) At(x, y int64) (*Chunk, error) {
	_, exists := world.data[util.Coords2i{X: x, Y: y}]

	// generate Chunk, if it doesn't exist yet
	if !exists {
		err := world.GenerateChunk(x, y)
		if err != nil {
			return nil, err
		}
	}

	return world.data[util.Coords2i{X: x, Y: y}], nil
}
