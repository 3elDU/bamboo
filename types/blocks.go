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
	CampfireBlock
	BerryBushBlock
	SandWithStonesBlock
	SandWithClayBlock
	PitBlock
)

func NewBlock(id BlockType) Block {
	switch id {
	case EmptyBlock:
		return NewEmptyBlock()
	case StoneBlock:
		return NewStoneBlock()
	case WaterBlock:
		return NewWaterBlock()
	case SandBlock:
		return NewSandBlock()
	case GrassBlock:
		return NewGrassBlock()
	case SnowBlock:
		return NewSnowBlock()
	case ShortGrassBlock:
		return NewShortGrassBlock()
	case TallGrassBlock:
		return NewTallGrassBlock()
	case FlowersBlock:
		return NewFlowersBlock()
	case PineTreeBlock:
		return NewPineTreeBlock()
	case RedMushroomBlock:
		return NewRedMushroomBlock()
	case WhiteMushroomBlock:
		return NewWhiteMushroomBlock()
	case CaveEntranceBlock:
		return NewCaveEntranceBlock(uuid.New())
	case CaveWallBlock:
		return NewCaveWallBlock()
	case CaveFloorBlock:
		return NewCaveFloorBlock()
	case CaveExitBlock:
		return NewCaveExitBlock()
	case PineSaplingBlock:
		return NewPineSaplingBlock()
	case CampfireBlock:
		return NewCampfireBlock()
	case BerryBushBlock:
		return NewBerryBushBlock(0)
	case SandWithStonesBlock:
		return NewSandWithStonesBlock()
	case SandWithClayBlock:
		return NewSandWithClayBlock()
	case PitBlock:
		return NewPitBlock()
	}

	return NewEmptyBlock()
}

var (
	NewEmptyBlock          func() Block
	NewStoneBlock          func() Block
	NewWaterBlock          func() Block
	NewSandBlock           func() Block
	NewGrassBlock          func() Block
	NewSnowBlock           func() Block
	NewShortGrassBlock     func() Block
	NewTallGrassBlock      func() Block
	NewFlowersBlock        func() Block
	NewPineTreeBlock       func() Block
	NewRedMushroomBlock    func() Block
	NewWhiteMushroomBlock  func() Block
	NewCaveEntranceBlock   func(uuid uuid.UUID) Block
	NewCaveWallBlock       func() Block
	NewCaveFloorBlock      func() Block
	NewCaveExitBlock       func() Block
	NewPineSaplingBlock    func() Block
	NewCampfireBlock       func() Block
	NewBerryBushBlock      func(berries int) Block
	NewSandWithStonesBlock func() Block
	NewSandWithClayBlock   func() Block
	NewPitBlock            func() Block
)

type Block interface {
	Type() BlockType

	Coords() Vec2u
	SetCoords(coords Vec2u)
	ParentChunk() Chunk
	SetParentChunk(chunk Chunk)

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
	// recursiveRedraw means that a block can trigger a redraw of other blocks/chunks
	// if it is set to false, the block shouldn't attempt to trigger redraw of other chunks/blocks
	Render(world World, screen *ebiten.Image, pos Vec2f, recursiveRedraw bool)
	TextureName() string
}

// A block that reacts to player colliding with it
type CollisionReactiveBlock interface {
	Block
	// Called when player collides(touches) with the block
	Collide(world World, playerPosition Vec2f)
}

// A block that the player can interact with
type InteractiveBlock interface {
	Interact()
}

// A block that can be broken. Does not necessarily mean that block drops an item.
type BreakableBlock interface {
	Block
	ToolRequiredToBreak() ToolFamily
	ToolStrengthRequired() ToolStrength
	Break()
}

type ICampfireBlock interface {
	AddPiece(item IBurnableItem) bool
	LightUp() bool
	IsLitUp() bool
}

// A generic crop block that can run out of water, and can be watered
type ICropBlock interface {
	NeedsWatering() bool
	AddWater()
}
