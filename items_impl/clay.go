package items_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewClayItem = NewClayItem
}

type ClayItem struct {
	baseItem
}

func NewClayItem() types.Item {
	return &ClayItem{
		baseItem: baseItem{
			id: types.ClayItem,
		},
	}
}

func (item *ClayItem) Name() string {
	return "Clay"
}
func (item *ClayItem) Description() string {
	return ""
}

func (item *ClayItem) Texture() *ebiten.Image {
	return assets.Texture("clay").Texture()
}

func (item *ClayItem) Family() types.ToolFamily {
	return types.ToolFamilyNone
}
func (item *ClayItem) Strength() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (item *ClayItem) Use(pos types.Vec2u) {
	// Clay can be placed back on sand
	if types.GetCurrentWorld().BlockAt(pos.X, pos.Y).Type() == types.SandBlock {
		types.GetCurrentWorld().SetBlock(
			pos.X, pos.Y,
			types.NewSandWithClayBlock(),
		)
		types.GetInventory().RemoveItemByType(item.Type(), 1)
	}
}

func (item *ClayItem) State() interface{} {
	return item.baseItem.State()
}
func (item *ClayItem) LoadState(state interface{}) {
	item.baseItem.LoadState(state)
}
