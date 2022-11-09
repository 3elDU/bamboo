package world

import (
	"log"
	"math/rand"

	"github.com/3elDU/bamboo/util"
	"github.com/aquilax/go-perlin"
	"github.com/google/uuid"
)

type World struct {
	// Separate perlin noise generators for each layer
	bottomGenerator *perlin.Perlin
	groundGenerator *perlin.Perlin
	topGenerator    *perlin.Perlin

	Metadata WorldSave

	// keeping registry of chunks that were loaded
	// in order from oldest to newest
	loadedChunks []util.Coords2i

	// keys there are Chunk coordinates.
	// so, actual Chunk coordinates are x*16 and y*16
	data map[util.Coords2i]*Chunk
}

// Creates a new world, using given name and seed
func NewWorld(name string, uuid uuid.UUID, seed int64) *World {
	log.Printf("NewWorld - name %v; seed %v", name, seed)

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

		Metadata: WorldSave{
			Name: name,
			UUID: uuid,
			Seed: seed,
		},

		loadedChunks: make([]util.Coords2i, 0),

		data: make(map[util.Coords2i]*Chunk),
	}
}

// Update - x and y are player coordinates
func (world *World) Update(_, _ float64) {

}

// At Returns a Chunk at given coordinates. Note that x and y are Chunk
// coordinates, not block coordinates
func (world *World) ChunkAt(x, y int64) (*Chunk, error) {
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

// Calculates chunk position from given world coordinates
// Acceps float64 so that negative coordinates will be handled properly
// Note that x and y are block coordinates
func (world *World) At(x, y float64) (*Chunk, error) {
	// HACK: handle negative coordinates properly
	if x < 0 {
		x -= 1
	}
	if y < 0 {
		y -= 1
	}
	var (
		cx = int64(x / 16)
		cy = int64(y / 16)
	)

	_, exists := world.data[util.Coords2i{X: cx, Y: cy}]

	if !exists {
		// try to load the chunk from disk first
		if chunk := LoadChunk(world.Metadata.UUID, cx, cy); chunk != nil {
			world.data[util.Coords2i{X: cx, Y: cy}] = chunk
		} else {
			chunk := NewChunk(cx, cy)
			err := chunk.Generate(world.bottomGenerator, world.groundGenerator, world.topGenerator)
			if err != nil {
				return nil, err
			}
			world.data[util.Coords2i{X: cx, Y: cy}] = chunk
		}

		world.loadedChunks = append(world.loadedChunks, util.Coords2i{X: cx, Y: cy})
	}

	return world.data[util.Coords2i{X: cx, Y: cy}], nil
}

func (world World) Seed() int64 {
	return world.Metadata.Seed
}
