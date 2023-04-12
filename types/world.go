package types

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type World interface {
	BlockAt(bx uint64, by uint64) Block
	CheckNeighbors(cx uint64, cy uint64) bool
	ChunkAt(cx uint64, cy uint64) Chunk
	ChunkAtB(bx uint64, by uint64) Chunk
	ChunkExists(cx uint64, cy uint64) bool
	GetNeighbors(cx uint64, cy uint64) []Chunk
	Render(screen *ebiten.Image, playerX float64, playerY float64, scaling float64)
	Save()
	Seed() int64
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
}
