package types

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type BlockType int

type Block interface {
	Coords() Vec2u
	SetCoords(coords Vec2u)
	ParentChunk() Chunk
	SetParentChunk(chunk Chunk)
	Type() BlockType

	Update(world World)

	State() interface{}
	// LoadState panicks on error
	LoadState(interface{})
}

type CollidableBlock interface {
	Block
	Collidable() bool
	// Collision points go in order: top-left, top-right, bottom-left, bottom-right
	CollisionPoints() [4]Vec2f
	PlayerSpeed() float64
}

type DrawableBlock interface {
	Block
	Render(world World, screen *ebiten.Image, pos Vec2f)
	TextureName() string
}

type InteractiveBlock interface {
	Block
	Interact(world World, playerPosition Vec2f)
}
