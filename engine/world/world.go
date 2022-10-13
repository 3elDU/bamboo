package world

import (
	"fmt"

	"github.com/3elDU/bamboo/config"
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

func (world *World) gen(x, y float64) Block {
	// returns a value from 0 to 2
	h := world.generator.Noise2D(x/config.PerlinNoiseScaleFactor, y/config.PerlinNoiseScaleFactor) + 1

	switch {
	case h <= 1: // Water
		return NewWaterBlock()
	case h <= 1.1: // Sand
		return NewSandBlock()
	case h <= 1.45: // Grass
		return NewGrassBlock()
	case h <= 1.65: // Stone
		return NewStoneBlock(h)
	default: // Snow
		return NewSnowBlock()
	}
}

func (world *World) GenerateChunk(cx, cy int64) error {
	fmt.Printf("Generating Chunk %v, %v\n", cx, cy)

	chunk := NewChunk(cx, cy)
	world.data[util.Coords2i{X: cx, Y: cy}] = chunk

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			b := world.gen(float64(int(cx)*16+x), float64(int(cy)*16+y))
			err := chunk.SetBlock(x, y, b)
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

func (world *World) Seed() int64 {
	return world.seed
}
