package blocks_impl

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(FlowersState{})
	types.NewFlowersBlock = NewFlowersBlock
}

type FlowersState struct {
	BaseBlockState
	TexturedBlockState
}

type FlowersBlock struct {
	baseBlock
	texturedBlock
}

func NewFlowersBlock() types.Block {
	return &FlowersBlock{
		baseBlock: baseBlock{
			blockType: types.FlowersBlock,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("flowers"),
			rotation: 0,
		},
	}
}

func (b *FlowersBlock) State() interface{} {
	return FlowersState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *FlowersBlock) LoadState(s interface{}) {
	state := s.(FlowersState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
