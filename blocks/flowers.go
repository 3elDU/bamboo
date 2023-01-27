package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(FlowersState{})
}

type FlowersState struct {
	BaseBlockState
	TexturedBlockState
}

type flowers struct {
	baseBlock
	texturedBlock
}

func NewFlowersBlock() *flowers {
	return &flowers{
		baseBlock: baseBlock{
			blockType: Flowers,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("flowers"),
			rotation: 0,
		},
	}
}

func (b flowers) State() interface{} {
	return FlowersState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *flowers) LoadState(s interface{}) {
	state := s.(FlowersState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
