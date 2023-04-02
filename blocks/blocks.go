package blocks

import (
	"github.com/3elDU/bamboo/types"
)

const (
	Empty types.BlockType = iota
	Stone
	Water
	Sand
	Grass
	Snow
	ShortGrass
	TallGrass
	Flowers
	PineTree
	RedMushroom
	WhiteMushroom
)

// GetBlockByID returns an empty block
func GetBlockByID(id types.BlockType) types.Block {
	switch id {
	case Empty:
		return NewEmptyBlock()
	case Stone:
		return NewStoneBlock()
	case Water:
		return NewWaterBlock()
	case Sand:
		return NewSandBlock(false)
	case Grass:
		return NewGrassBlock()
	case Snow:
		return NewSnowBlock()
	case ShortGrass:
		return NewShortGrassBlock()
	case TallGrass:
		return NewTallGrassBlock()
	case Flowers:
		return NewFlowersBlock()
	case PineTree:
		return NewPineTreeBlock()
	case RedMushroom:
		return NewRedMushroomBlock()
	case WhiteMushroom:
		return NewWhiteMushroomBlock()
	}

	return NewEmptyBlock()
}
