package world

import (
	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/worldgen"
	"log"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
)

type World struct {
	generator   types.WorldGenerator
	saverLoader *SaverLoader

	metadata types.Save

	chunks map[types.Vec2u]*Chunk
}

// Creates a new world, using given name and seed
func NewWorld(metadata types.Save) *World {
	log.Printf("NewWorld - name %v; seed %v", metadata.Name, metadata.Seed)

	generator := worldgen.NewWorldgenForType(metadata.Seed, metadata.WorldType)
	go generator.Run()

	saverLoader := NewWorldSaverLoader(metadata)
	go saverLoader.Run()

	return &World{
		generator:   generator,
		saverLoader: saverLoader,

		metadata: metadata,

		chunks: make(map[types.Vec2u]*Chunk),
	}
}

// Update - x and y are player coordinates
func (world *World) Update() {
	// receive newly generated chunks from world generator
	chunks := world.generator.Receive()
	for _, chunk := range chunks {
		world.chunks[chunk.Coords()] = chunk.(*Chunk)
		// Request redraw of each neighbor
		for _, neighbor := range world.GetNeighbors(chunk.Coords().X, chunk.Coords().Y) {
			neighbor.TriggerRedraw()
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
		for coords, chunk := range world.chunks {
			if scene_manager.Ticks()-chunk.lastAccessed > config.ChunkUnloadDelay {
				world.saverLoader.Save(chunk)
				delete(world.chunks, coords)
			}
		}
	}

	// Update all currently loaded chunks
	for _, chunk := range world.chunks {
		chunk.Update(world)
	}
}

func (world *World) ChunkAt(cx, cy uint64) types.Chunk {
	return world.ChunkAtB(cx*16, cy*16)
}

func (world *World) ChunkAtB(bx, by uint64) types.Chunk {
	cx := bx / 16
	cy := by / 16
	chunkCoordinates := types.Vec2u{X: cx, Y: cy}

	_, exists := world.chunks[chunkCoordinates]

	if !exists {
		// try to load the chunk from disk first
		if ChunkExistsOnDisk(world.metadata, cx, cy) {
			// request chunk loading
			world.saverLoader.Load(cx, cy)
		} else {
			// request chunk generation
			chunk := NewChunk(cx, cy)
			world.generator.Generate(chunk)
		}
		dummyChunk := NewChunk(cx, cy)
		world.generator.GenerateDummy(dummyChunk)
		world.chunks[chunkCoordinates] = dummyChunk
	}

	world.chunks[chunkCoordinates].lastAccessed = scene_manager.Ticks()
	return world.chunks[chunkCoordinates]
}

// There is no B suffix, because it's trivial that this function accepts block coordinates
func (world *World) BlockAt(bx, by uint64) types.Block {
	cx, cy := bx/16, by/16

	chunk, exists := world.chunks[types.Vec2u{X: cx, Y: cy}]
	if !exists {
		return blocks.NewEmptyBlock()
	}

	return chunk.At(uint(bx%16), uint(by%16))
}

func (world *World) SetBlock(bx, by uint64, block types.Block) {
	cx, cy := bx/16, by/16

	if !world.ChunkExists(cx, cy) {
		// generate a chunk immediately, if it doesn't exist
		c := NewChunk(cx, cy)
		world.generator.GenerateImmediately(c)
		world.chunks[types.Vec2u{X: cx, Y: cy}] = c
	}

	world.chunks[types.Vec2u{X: cx, Y: cy}].SetBlock(uint(bx%16), uint(by%16), block)
}

func (world *World) ChunkExists(cx, cy uint64) bool {
	_, exists := world.chunks[types.Vec2u{X: cx, Y: cy}]
	return exists
}

func (world *World) GetNeighbors(cx, cy uint64) []types.Chunk {
	sides := [4]types.Vec2u{
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

func (world *World) Seed() int64 {
	return world.metadata.Seed
}

func (world *World) Metadata() types.Save {
	return world.metadata
}

func (world *World) Generator() types.WorldGenerator {
	return world.generator
}
