package items_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewRawIronItem = NewRawIronItem
}

type RawIronItem struct {
	baseItem
}

func NewRawIronItem() types.Item {
	return &RawIronItem{
		baseItem{id: types.RawIronItem},
	}
}

func (item *RawIronItem) Name() string {
	return "Raw iron"
}
func (item *RawIronItem) Description() string {
	return ""
}
func (item *RawIronItem) Texture() *ebiten.Image {
	return assets.Texture("raw_iron").Texture()
}

func (item *RawIronItem) SmeltingEnergyRequired() float64 {
	return 5.0
}
func (item *RawIronItem) Smelt() types.Item {
	return types.NewIronIngotItem()
}
