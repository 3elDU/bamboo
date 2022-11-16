package world

import (
	"log"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/scene_manager"
	"github.com/3elDU/bamboo/util"
	"github.com/google/uuid"
)

type World struct {
	generator   *WorldGenerator
	saverLoader *WorldSaverLoader

	Metadata WorldSave

	// keys there are Chunk coordinates.
	// so, actual Chunk coordinates are x*16 and y*16
	chunks map[util.Coords2i]*Chunk
}

// Creates a new world, using given name and seed
func NewWorld(name string, uuid uuid.UUID, seed int64) *World {
	log.Printf("NewWorld - name %v; seed %v", name, seed)

	generator := NewWorldGenerator(seed)
	go generator.Run()

	metadata := WorldSave{Name: name, UUID: uuid, Seed: seed}

	saverLoader := NewWorldSaverLoader(metadata)
	go saverLoader.Run()

	return &World{
		generator:   generator,
		saverLoader: saverLoader,

		Metadata: metadata,

		chunks: make(map[util.Coords2i]*Chunk),
	}
}

// Update - x and y are player coordinates
func (world *World) Update(_, _ float64) {
	// receive newly generated chunks from world generator
	for {
		if chunk := world.generator.Receive(); chunk != nil {
			world.chunks[chunk.Coords()] = chunk
		} else {
			break
		}
	}

	// receive newly loaded chunks
	for {
		if chunk := world.saverLoader.Receive(); chunk != nil {
			world.chunks[chunk.Coords()] = chunk
		} else {
			break
		}
	}

	// each 30 ticks ( half a second ) check for chunks,
	// that weren't accessed ( neither read, nor write ) for specified amount of ticks
	// ( check config.go )
	if scene_manager.Ticks()%30 == 0 {
		chunksUnloaded := 0

		for coords, chunk := range world.chunks {
			if scene_manager.Ticks()-chunk.lastAccessed > config.ChunkUnloadDelay {
				world.saverLoader.Save(chunk)
				delete(world.chunks, coords)
				chunksUnloaded++
			}
		}

		if chunksUnloaded > 0 {
			log.Printf("World.Update() - Unloaded %v chunks from memory; currently loaded - %v", chunksUnloaded, len(world.chunks))
		}
	}
}

// At Returns a Chunk at given coordinates. Note that x and y are Chunk
// coordinates, not block coordinates
func (world *World) ChunkAt(x, y int64) (*Chunk, error) {
	return world.At(float64(x)*16, float64(y)*16)
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
	chunkCoordinates := util.Coords2i{X: cx, Y: cy}

	_, exists := world.chunks[chunkCoordinates]

	if !exists {
		// try to load the chunk from disk first
		if ChunkExistsOnDisk(world.Metadata.UUID, cx, cy) {
			// request chunk loading
			world.saverLoader.Load(cx, cy)
		} else {
			// request chunk generation
			world.generator.Generate(cx, cy)
		}

		world.chunks[chunkCoordinates] = NewDummyChunk(cx, cy)
		return world.chunks[chunkCoordinates], nil
	}

	world.chunks[chunkCoordinates].lastAccessed = scene_manager.Ticks()
	return world.chunks[chunkCoordinates], nil
}

func (world World) Seed() int64 {
	return world.Metadata.Seed
}
