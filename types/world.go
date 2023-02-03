package types

import "github.com/hajimehoshi/ebiten/v2"

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
