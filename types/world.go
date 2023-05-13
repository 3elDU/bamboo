package types

import (
	"github.com/3elDU/bamboo/world_type"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

var currentWorld World

// Sets a reference to currently loaded world
func SetCurrentWorld(world World) {
	currentWorld = world
}

// Returns reference to currently loaded World instance
func GetCurrentWorld() World {
	return currentWorld
}

type World interface {
	// There is no B suffix, because it's trivial that this function accepts block coordinates
	BlockAt(bx uint64, by uint64) Block
	SetBlock(bx, by uint64, block Block)
	// Checks neighbors of the given chunk
	// Returns false if at least one of them doesn't exist
	// Automatically requests generation of neighbors
	CheckNeighbors(cx uint64, cy uint64) bool
	// Checks if each of 8 neighboring blocks type is in the list of allowed types
	CheckBlockNeighbors(bx, by uint64, allowedTypes []BlockType) bool
	// Checks if at least one of neighboring blocks type is in the list
	BlockNeighboringWith(bx, by uint64, types []BlockType) bool
	ChunkAt(cx uint64, cy uint64) Chunk
	ChunkAtB(bx uint64, by uint64) Chunk
	ChunkExists(cx uint64, cy uint64) bool
	GetNeighbors(cx uint64, cy uint64) []Chunk
	// Returns world generator associated with this world
	Generator() WorldGenerator
	Metadata() Save
	Render(screen *ebiten.Image, playerX float64, playerY float64, scaling float64)
	Save()
	Seed() int64
	// Returned size is in chunks
	Size() Vec2u
	Update()
}

// Structure with metadata, representing a world save
type Save struct {
	Name string // world name as from the user
	// BaseUUID is a base uuid for all worlds, as well as name of the base save folder.
	// All subsequent worlds will be created in their own directories, under the base directory.
	BaseUUID uuid.UUID
	// ID of the current world. There can be many worlds in a save.
	UUID uuid.UUID
	Seed int64
	// Size of the world in blocks
	Size      Vec2u
	WorldType world_type.WorldType
}
