package items_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewClayPickaxeItem = NewClayPickaxeItem
}

type ClayPickaxeItem struct {
	baseItem
}

func NewClayPickaxeItem() types.Item {
	return &ClayPickaxeItem{
		baseItem{id: types.ClayPickaxeItem},
	}
}

func (pickaxe *ClayPickaxeItem) Name() string {
	return "Clay pickaxe"
}
func (pickaxe *ClayPickaxeItem) Description() string {
	return ""
}
func (pickaxe *ClayPickaxeItem) Texture() *ebiten.Image {
	return assets.Texture("clay_pickaxe").Texture()
}

func (pickaxe *ClayPickaxeItem) ToolFamily() types.ToolFamily {
	return types.ToolFamilyPickaxe
}
func (pickaxe *ClayPickaxeItem) ToolStrength() types.ToolStrength {
	return types.ToolStrengthClay
}
func (pickaxe *ClayPickaxeItem) UseTool(_ types.Vec2u) {

}
