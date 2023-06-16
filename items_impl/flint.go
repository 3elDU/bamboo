package items_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
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
	return asset_loader.Texture("flint").Texture()
}

func (flint *FlintItem) Use(pos types.Vec2u) {
	block := types.GetCurrentWorld().BlockAt(pos.X, pos.Y)
	switch block.Type() {
	case types.CampfireBlock:
		if block.Type() != types.CampfireBlock {
			return
		}

		campfire := block.(types.CampfireBlockI)
		litUp := campfire.LightUp()

		if litUp {
			types.GetInventory().RemoveItem(types.ItemSlot{
				Item:     flint,
				Quantity: 1,
			})
		}
	}
}
