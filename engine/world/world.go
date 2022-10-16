package world

import (
	"math/rand"

	"github.com/3elDU/bamboo/util"
	"github.com/aquilax/go-perlin"
)

type World struct {
	// Separate perlin noise generators for each layer
	bottomGenerator *perlin.Perlin
	groundGenerator *perlin.Perlin
	topGenerator    *perlin.Perlin

	mapSeed int64

	// keys there are Chunk coordinates.
	// so, actual Chunk coordinates are x*16 and y*16
	data map[util.Coords2i]*Chunk
}

func NewWorld(seed int64) *World {
	// make a random generator using global world seed
	world := rand.New(rand.NewSource(seed))

	// and generate perlin noise seeds, using it
	var (
		bottomSeed = world.Int63()
		groundSeed = world.Int63()
		topSeed    = world.Int63()
	)

	return &World{
		bottomGenerator: perlin.NewPerlin(2, 2, 16, bottomSeed),
		groundGenerator: perlin.NewPerlin(2, 2, 16, groundSeed),
		topGenerator:    perlin.NewPerlin(2, 2, 16, topSeed),

		mapSeed: seed,

		data: make(map[util.Coords2i]*Chunk),
	}
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
		chunk := NewChunk(x, y)
		err := chunk.Generate(world.bottomGenerator, world.groundGenerator, world.topGenerator)
		if err != nil {
			return nil, err
		}
		world.data[util.Coords2i{X: x, Y: y}] = chunk
	}

	return world.data[util.Coords2i{X: x, Y: y}], nil
}

func (world *World) Seed() int64 {
	return world.mapSeed
}
