package items_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(FlintItemState{})
	types.NewFlintItem = NewFlintItem
}

type FlintItemState struct {
	BaseItemState
}

type FlintItem struct {
	baseItem
}

func NewFlintItem() types.Item {
	return &FlintItem{
		baseItem: baseItem{
			id: types.FlintItem,
		},
	}
}

func (flint *FlintItem) Name() string {
	return "Flint"
}

func (flint *FlintItem) Description() string {
	return ""
}

func (flint *FlintItem) Hash() uint64 {
	return uint64(flint.id)
}

func (flint *FlintItem) Texture() *ebiten.Image {
	return assets.Texture("flint").Texture()
}

func (item *FlintItem) ToolFamily() types.ToolFamily {
	return types.ToolFamilyNone
}
func (item *FlintItem) ToolStrength() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (flint *FlintItem) UseTool(pos types.Vec2u) {
	block := types.GetCurrentWorld().BlockAt(pos.X, pos.Y)
	switch block.Type() {
	case types.CampfireBlock:
		if block.Type() != types.CampfireBlock {
			return
		}

		campfire := block.(types.ICampfireBlock)
		if campfire.LightUp() {
			types.GetPlayerInventory().RemoveItemByType(flint.Type(), 1)
		}
	case types.SandBlock:
		// Flint can be placed back on sand
		types.GetCurrentWorld().SetBlock(pos.X, pos.Y, types.NewSandWithStonesBlock())
		types.GetPlayerInventory().RemoveItemByType(flint.Type(), 1)
	}
}
