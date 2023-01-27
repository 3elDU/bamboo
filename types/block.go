/*
	Declarations of basic block types
*/

package types

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type BlockType int

type Block interface {
	Coords() Coords2u
	SetCoords(coords Coords2u)
	ParentChunk() Chunk
	SetParentChunk(chunk Chunk)
	Type() BlockType

	Update(world World)

	State() interface{}
	// panicks on error
	LoadState(interface{})
}

type CollidableBlock interface {
	Block
	Collidable() bool
	// Collision points go in order: top-left, top-right, bottom-left, bottom-right
	CollisionPoints() [4]Coords2f
	PlayerSpeed() float64
}

type DrawableBlock interface {
	Block
	Render(world World, screen *ebiten.Image, pos Coords2f)
	TextureName() string
}
