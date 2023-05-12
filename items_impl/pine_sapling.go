package items_impl

import (
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewPineSaplingItem = NewPineSaplingItem
}

type PineSaplingItem struct {
	baseItem
	Tex types.Texture
}

func NewPineSaplingItem() types.Item {
	return &PineSaplingItem{
		baseItem: baseItem{
			id: types.PineSaplingItem,
		},
		Tex: asset_loader.Texture("sapling_item"),
	}
}

func (item *PineSaplingItem) Texture() *ebiten.Image {
	return item.Tex.Texture()
}

func (item *PineSaplingItem) Hash() uint64 {
	return uint64(item.id)
}

func (item *PineSaplingItem) Use(pos types.Vec2u) {
	// sapling can only be planted on empty grass
	if types.GetCurrentWorld().BlockAt(pos.X, pos.Y).Type() != types.GrassBlock {
		return
	}
	types.GetCurrentWorld().SetBlock(pos.X, pos.Y, types.NewPineSaplingBlock())
	types.GetInventory().RemoveItem(types.ItemSlot{
		Item:     item,
		Quantity: 1,
	})
}
