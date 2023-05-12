package blocks_impl

import (
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
)

func init() {
	types.NewBlock = newBlockByID
}

func newBlockByID(id types.BlockType) types.Block {
	switch id {
	case types.EmptyBlock:
		return NewEmptyBlock()
	case types.StoneBlock:
		return NewStoneBlock()
	case types.WaterBlock:
		return NewWaterBlock()
	case types.SandBlock:
		return NewSandBlock(false)
	case types.GrassBlock:
		return NewGrassBlock()
	case types.SnowBlock:
		return NewSnowBlock()
	case types.ShortGrassBlock:
		return NewShortGrassBlock()
	case types.TallGrassBlock:
		return NewTallGrassBlock()
	case types.FlowersBlock:
		return NewFlowersBlock()
	case types.PineTreeBlock:
		return NewPineTreeBlock()
	case types.RedMushroomBlock:
		return NewRedMushroomBlock()
	case types.WhiteMushroomBlock:
		return NewWhiteMushroomBlock()
	case types.CaveEntranceBlock:
		return NewCaveEntranceBlock(uuid.New())
	case types.CaveWallBlock:
		return NewCaveWallBlock()
	case types.CaveFloorBlock:
		return NewCaveFloorBlock()
	case types.CaveExitBlock:
		return NewCaveExitBlock()
	case types.PineSaplingBlock:
		return NewPineSaplingBlock()
	}

	return NewEmptyBlock()
}
