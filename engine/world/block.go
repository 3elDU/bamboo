package world

import (
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Block interface {
	Coords() util.Coords2i
	SetCoords(coords util.Coords2i)
	ParentChunk() *Chunk
	SetParentChunk(chunk *Chunk)

	// Whether the player should collide with the block
	Collidable() bool

	Update()
	Render(screen *ebiten.Image, pos util.Coords2f)
}
