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
		return NewSandBlock(false)
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
	}

	return NewEmptyBlock()
}

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
	NewCampfireBlock      func() Block
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

// A block that can be broken. Does not necessarily mean that block drops an item.
type BreakableBlock interface {
	Block
	Break()
}

type CampfireBlockI interface {
	AddPiece(item BurnableItem)
	LightUp() bool
}
