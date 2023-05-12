package types

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type BlockType int

const (
	EmptyBlock BlockType = iota
	StoneBlock
	WaterBlock
	SandBlock
	GrassBlock
	SnowBlock
	ShortGrassBlock
	TallGrassBlock
	FlowersBlock
	PineTreeBlock
	RedMushroomBlock
	WhiteMushroomBlock
	CaveEntranceBlock
	CaveWallBlock
	CaveFloorBlock
	CaveExitBlock
	PineSaplingBlock
)

var NewBlock func(id BlockType) Block
var (
	NewEmptyBlock         func() Block
	NewStoneBlock         func() Block
	NewWaterBlock         func() Block
	NewSandBlock          func(stones bool) Block
	NewGrassBlock         func() Block
	NewSnowBlock          func() Block
	NewShortGrassBlock    func() Block
	NewTallGrassBlock     func() Block
	NewFlowersBlock       func() Block
	NewPineTreeBlock      func() Block
	NewRedMushroomBlock   func() Block
	NewWhiteMushroomBlock func() Block
	NewCaveEntranceBlock  func(uuid uuid.UUID) Block
	NewCaveWallBlock      func() Block
	NewCaveFloorBlock     func() Block
	NewCaveExitBlock      func() Block
	NewPineSaplingBlock   func() Block
)

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

// A block that player can collide with
type CollidableBlock interface {
	Block
	Collidable() bool
	// Collision points go in order: top-left, top-right, bottom-left, bottom-right
	CollisionPoints() [4]Vec2f
	PlayerSpeed() float64
}

// A block that can be rendered onto the screen
type DrawableBlock interface {
	Block
	Render(world World, screen *ebiten.Image, pos Vec2f)
	TextureName() string
}

// A block that player can interact with
type InteractiveBlock interface {
	Block
	Interact(world World, playerPosition Vec2f)
}

// A block that can be broken, and drops an item
type BreakableBlock interface {
	Block
	Break()
}
