package items_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewClayShovelItem = NewClayShovelItem
}

type ClayShovelItem struct {
	baseItem
}

func NewClayShovelItem() types.Item {
	return &ClayShovelItem{
		baseItem: baseItem{
			id: types.ClayShovelItem,
		},
	}
}

func (shovel *ClayShovelItem) Stackable() bool {
	return false
}

func (shovel *ClayShovelItem) Name() string {
	return "Clay shovel"
}
func (shovel *ClayShovelItem) Description() string {
	return ""
}

func (shovel *ClayShovelItem) Texture() *ebiten.Image {
	return assets.Texture("clay_shovel").Texture()
}

func (shovel *ClayShovelItem) Family() types.ToolFamily {
	return types.ToolFamilyShovel
}
func (shovel *ClayShovelItem) Strength() types.ToolStrength {
	return types.ToolStrengthClay
}
func (shovel *ClayShovelItem) Use(pos types.Vec2u) {

}
