package blocks_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	types.NewSandWithStonesBlock = NewSandWithStonesBlock
}

type SandWithStonesBlock struct {
	baseBlock
	texturedBlock
	collidableBlock
}

func NewSandWithStonesBlock() types.Block {
	return &SandWithStonesBlock{
		baseBlock: baseBlock{
			blockType: types.SandWithStonesBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture("sand_stones"),
		},
		collidableBlock: collidableBlock{
			collidable:  false,
			playerSpeed: 0.8,
		},
	}
}

func (b *SandWithStonesBlock) ToolRequiredToBreak() types.ToolFamily {
	return types.ToolFamilyNone
}
func (b *SandWithStonesBlock) ToolStrengthRequired() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (b *SandWithStonesBlock) Break() {
	if types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewFlintItem(),
		Quantity: 1,
	}) {
		types.GetCurrentWorld().SetBlock(uint64(b.x), uint64(b.y), types.NewSandBlock())
	}
}

func (b *SandWithStonesBlock) State() interface{} {
	return b.baseBlock.State()
}
func (b *SandWithStonesBlock) LoadState(state interface{}) {
	b.baseBlock.LoadState(state)
}
