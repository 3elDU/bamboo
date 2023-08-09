package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(ShortGrassState{})
	types.NewShortGrassBlock = NewShortGrassBlock
}

type ShortGrassState struct {
	BaseBlockState
}

type ShortGrassBlock struct {
	connectedBlock
}

func NewShortGrassBlock() types.Block {
	return &ShortGrassBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.ShortGrassBlock,
			},
			tex: assets.ConnectedTexture("short_grass", true, true, true, true),
			connectsTo: []types.BlockType{
				types.ShortGrassBlock,
				types.FlowersBlock,
				types.RedMushroomBlock, types.WhiteMushroomBlock,
			},
		},
	}
}

func (b *ShortGrassBlock) ToolRequiredToBreak() types.ToolFamily {
	return types.ToolFamilyNone
}
func (b *ShortGrassBlock) ToolStrengthRequired() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (b *ShortGrassBlock) Break() {
	types.GetCurrentWorld().SetBlock(uint64(b.x), uint64(b.y), types.NewGrassBlock())
}

func (b *ShortGrassBlock) State() interface{} {
	return ShortGrassState{
		BaseBlockState: b.baseBlock.State().(BaseBlockState),
	}
}

func (b *ShortGrassBlock) LoadState(s interface{}) {
	state := s.(ShortGrassState)
	b.baseBlock.LoadState(state.BaseBlockState)
}
