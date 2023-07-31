package blocks_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	types.NewSandWithClayBlock = NewSandWithClayBlock
}

type SandWithClayBlock struct {
	baseBlock
	texturedBlock
	collidableBlock
}

func NewSandWithClayBlock() types.Block {
	return &SandWithClayBlock{
		baseBlock: baseBlock{
			blockType: types.SandWithClayBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture("sand_clay"),
		},
		collidableBlock: collidableBlock{
			collidable:  false,
			playerSpeed: 0.8,
		},
	}
}

func (b *SandWithClayBlock) Break() {
	if types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewClayItem(),
		Quantity: 1,
	}) {
		types.GetCurrentWorld().SetBlock(uint64(b.x), uint64(b.y), types.NewSandBlock())
	}
}

func (b *SandWithClayBlock) State() interface{} {
	return b.baseBlock.State()
}
func (b *SandWithClayBlock) LoadState(state interface{}) {
	b.baseBlock.LoadState(state)
}
