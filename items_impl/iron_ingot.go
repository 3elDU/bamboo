package items_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewIronIngotItem = NewIronIngotItem
}

type IronIngotItem struct {
	baseItem
}

func NewIronIngotItem() types.Item {
	return &IronIngotItem{
		baseItem{id: types.IronIngotItem},
	}
}

func (item *IronIngotItem) Name() string {
	return "Iron ingot"
}
func (item *IronIngotItem) Description() string {
	return ""
}
func (item *IronIngotItem) Texture() *ebiten.Image {
	return assets.Texture("iron_ingot").Texture()
}
