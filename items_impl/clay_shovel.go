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

func (shovel *ClayShovelItem) ToolFamily() types.ToolFamily {
	return types.ToolFamilyShovel
}
func (shovel *ClayShovelItem) ToolStrength() types.ToolStrength {
	return types.ToolStrengthClay
}
func (shovel *ClayShovelItem) UseTool(pos types.Vec2u) {
	types.GetCurrentWorld().SetBlock(pos.X, pos.Y, types.NewPitBlock())
}
