package world

import (
	"log"

	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
)

type World struct {
	generator   *WorldGenerator
	saverLoader *WorldSaverLoader

	Metadata WorldSave

	// keys there are Chunk coordinates.
	// so, actual Chunk coordinates are x*16 and y*16
	chunks map[types.Coords2u]*Chunk
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

		chunks: make(map[types.Coords2u]*Chunk),
	}
}

// Update - x and y are player coordinates
func (world *World) Update() {
	// receive newly generated chunks from world generator
	for {
		if chunk := world.generator.Receive(); chunk != nil {
			log.Printf("world.Update() - received chunk %v, %v", chunk.Coords().X, chunk.Coords().Y)
			world.chunks[chunk.Coords()] = chunk
			// Request redraw of each neighbor
			for _, neighbor := range world.GetNeighbors(chunk.Coords().X, chunk.Coords().Y) {
				neighbor.TriggerRedraw()
			}
		} else {
			break
		}
	}

	// receive newly loaded chunks
	for {
		if chunk := world.saverLoader.Receive(); chunk != nil {
			world.chunks[chunk.Coords()] = chunk
			// Request redraw of each neighbor
			for _, neighbor := range world.GetNeighbors(chunk.Coords().X, chunk.Coords().Y) {
				neighbor.TriggerRedraw()
			}
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

	// updateStart := time.Now()
	// Update all currently loaded chunks
	for _, chunk := range world.chunks {
		chunk.Update(world)
	}
	// updateEnd := time.Now()
	// log.Printf("World.Update() - chunk update took %v; loaded chunks - %v", updateEnd.Sub(updateStart).String(), len(world.chunks))
}

func (world *World) ChunkAt(cx, cy uint64) types.Chunk {
	return world.ChunkAtB(cx*16, cy*16)
}

func (world *World) ChunkAtB(bx, by uint64) types.Chunk {
	cx := bx / 16
	cy := by / 16
	chunkCoordinates := types.Coords2u{X: cx, Y: cy}

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

		// world.chunks[chunkCoordinates] = NewDummyChunk(cx, cy)
		// return world.chunks[chunkCoordinates]
		return NewDummyChunk(cx, cy)
	}

	world.chunks[chunkCoordinates].lastAccessed = scene_manager.Ticks()
	return world.chunks[chunkCoordinates]
}

// There is no B suffix, because it's trivial that this function accepts block coordinates
func (world *World) BlockAt(bx, by uint64) types.Block {
	cx, cy := bx/16, by/16

	chunk, exists := world.chunks[types.Coords2u{X: cx, Y: cy}]
	if !exists {
		return blocks.NewEmptyBlock()
	}

	return chunk.At(uint(bx%16), uint(by%16))
}

func (world *World) ChunkExists(cx, cy uint64) bool {
	_, exists := world.chunks[types.Coords2u{X: cx, Y: cy}]
	return exists
}

func (world *World) GetNeighbors(cx, cy uint64) []types.Chunk {
	sides := [4]types.Coords2u{
		{X: cx - 1, Y: cy}, // left
		{X: cx + 1, Y: cy}, // right
		{X: cx, Y: cy - 1}, // top
		{X: cx, Y: cy + 1}, // bottom
	}

	neighbors := make([]types.Chunk, 0)

	for _, side := range sides {
		if !world.ChunkExists(side.X, side.Y) {
			continue
		}
		neighbors = append(neighbors, world.ChunkAt(side.X, side.Y))
	}

	return neighbors
}

// Checks neighbors of the given chunk
// Returns false if at least one of them doesn't exist
// Automatically requests generation of neighbors
func (world *World) CheckNeighbors(cx, cy uint64) bool {
	if !world.ChunkExists(cx, cy) {
		// If the given chunk doesn't exist
		return false
	}

	return len(world.GetNeighbors(cx, cy)) == 4
}

func (world World) Seed() int64 {
	return world.Metadata.Seed
}
