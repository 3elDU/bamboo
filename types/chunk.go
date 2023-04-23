package types

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Chunk interface {
	// Returns a dummy block, in case of an error
	At(x uint, y uint) Block
	BlockCoords() Vec2u
	Coords() Vec2u
	Render(world World)
	Save(metadata Save)
	SetBlock(x uint, y uint, block Block)
	Update(world World)
	TriggerRedraw()
	Texture() *ebiten.Image
}
