package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(GrassBlockState{})
	types.NewGrassBlock = NewGrassBlock
}

type GrassBlockState struct {
	ConnectedBlockState
}

type GrassBlock struct {
	connectedBlock
}

func NewGrassBlock() types.Block {
	return &GrassBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.GrassBlock,
			},
			tex: assets.ConnectedTexture("grass", true, true, true, true),
			connectsTo: []types.BlockType{
				types.GrassBlock, types.ShortGrassBlock, types.TallGrassBlock, types.FlowersBlock, types.PineSaplingBlock, types.BerryBushBlock,
				types.PineTreeBlock,
				types.RedMushroomBlock, types.WhiteMushroomBlock,
				types.StoneBlock,
				types.CaveEntranceBlock,
				types.CampfireBlock,
			},
		},
	}
}

func (b *GrassBlock) State() interface{} {
	return GrassBlockState{
		ConnectedBlockState: b.connectedBlock.State().(ConnectedBlockState),
	}
}

func (b *GrassBlock) LoadState(s interface{}) {
	state := s.(GrassBlockState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
}
