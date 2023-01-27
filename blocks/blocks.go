/*
	Implementations for various block types
*/

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
	Short_Grass
	Tall_Grass
	Flowers
	PineTree
	RedMushroom
	WhiteMushroom
)

// Returns an empty interface
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
	case Short_Grass:
		return NewShortGrassBlock()
	case Tall_Grass:
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
