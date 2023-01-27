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

type snow struct {
	baseBlock
	texturedBlock
}

func NewSnowBlock() *snow {
	return &snow{
		baseBlock: baseBlock{
			blockType: Short_Grass,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("short_grass"),
			rotation: 0,
		},
	}
}

func (b snow) State() interface{} {
	return SnowState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *snow) LoadState(s interface{}) {
	state := s.(SnowState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
