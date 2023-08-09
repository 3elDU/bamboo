package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/assets"
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
			tex:      assets.Texture("flowers"),
			rotation: 0,
		},
	}
}

func (b *FlowersBlock) ToolRequiredToBreak() types.ToolFamily {
	return types.ToolFamilyNone
}
func (b *FlowersBlock) ToolStrengthRequired() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (b *FlowersBlock) Break() {
	types.GetCurrentWorld().SetBlock(uint64(b.x), uint64(b.y), types.NewGrassBlock())
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
