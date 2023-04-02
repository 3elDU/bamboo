package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(SnowState{})
}

type SnowState struct {
	BaseBlockState
	TexturedBlockState
}

type SnowBlock struct {
	baseBlock
	texturedBlock
}

func NewSnowBlock() *SnowBlock {
	return &SnowBlock{
		baseBlock: baseBlock{
			blockType: ShortGrass,
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
