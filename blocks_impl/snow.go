package blocks_impl

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(SnowState{})
	types.NewSnowBlock = NewSnowBlock
}

type SnowState struct {
	BaseBlockState
	TexturedBlockState
}

type SnowBlock struct {
	baseBlock
	texturedBlock
}

func NewSnowBlock() types.Block {
	return &SnowBlock{
		baseBlock: baseBlock{
			blockType: types.ShortGrassBlock,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("short_grass"),
			rotation: 0,
		},
	}
}

func (b *SnowBlock) State() interface{} {
	return SnowState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *SnowBlock) LoadState(s interface{}) {
	state := s.(SnowState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
