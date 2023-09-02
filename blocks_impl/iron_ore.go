package blocks_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	types.NewIronOreBlock = NewIronOreBlock
}

type IronOreBlock struct {
	baseBlock
	texturedBlock
	breakableBlock
	collidableBlock
}

func NewIronOreBlock() types.Block {
	return &IronOreBlock{
		baseBlock: baseBlock{
			blockType: types.IronOreBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture("iron_ore"),
		},
		breakableBlock: breakableBlock{
			toolRequiredToBreak:  types.ToolFamilyPickaxe,
			toolStrengthRequired: types.ToolStrengthWood,
		},
		collidableBlock: collidableBlock{
			collidable: true,
		},
	}
}

func (block *IronOreBlock) Break() {
	if types.GetPlayerInventory().AddItem(types.ItemSlot{Item: types.NewRawIronItem(), Quantity: 1}) {
		types.GetCurrentWorld().SetBlock(uint64(block.x), uint64(block.y), types.NewCaveFloorBlock(false))
	}
}

// Dummy methods for saving/loading state
func (block *IronOreBlock) LoadState(_ interface{}) {

}
func (block *IronOreBlock) State() interface{} {
	return nil
}
