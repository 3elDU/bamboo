package inventory

import (
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

const InventorySize = 5

type Inventory struct {
	Slots        [InventorySize]*types.ItemSlot
	SelectedSlot int
}

func NewInventory() *Inventory {
	inv := Inventory{}
	for i := range inv.Slots {
		inv.Slots[i] = new(types.ItemSlot)
		inv.Slots[i].Empty = true
	}

	return &inv
}

// Returns false if there is no space
func (inv *Inventory) AddItem(item types.Item) bool {
	for _, slot := range inv.Slots {
		if slot.AddItem(item) {
			return true
		}
	}
	return false
}

func (inv *Inventory) SelectSlot(slot int) {
	if slot >= InventorySize {
		slot = 0
	} else if slot < 0 {
		slot = InventorySize - 1
	}

	inv.SelectedSlot = slot
}

func (inv *Inventory) ItemInHand() types.Item {
	return inv.Slots[inv.SelectedSlot].Item
}

func (inv *Inventory) Render(screen *ebiten.Image) {
	inventoryTexture := asset_loader.Texture("inventory")
	w, h := inventoryTexture.ScaledSize()
	inventoryDrawOpts := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()

	// position of inventory texture on the screen
	ix := float64(sw)/2 - float64(w)/2 // horizontally centered
	iy := float64(sh) - float64(h)     // bottom of the screen

	inventoryDrawOpts.GeoM.Scale(float64(config.UIScaling), float64(config.UIScaling))
	inventoryDrawOpts.GeoM.Translate(ix, iy)

	screen.DrawImage(inventoryTexture.Texture(), inventoryDrawOpts)

	for i, slot := range inv.Slots {
		if slot.Empty {
			continue
		}

		itemTex := slot.Item.Texture()
		itemTexOpts := &ebiten.DrawImageOptions{}

		itemTexOpts.GeoM.Scale(float64(config.UIScaling), float64(config.UIScaling))
		itemTexOpts.GeoM.Translate(
			ix+4*float64(config.UIScaling)+(20*float64(i)*float64(config.UIScaling)),
			iy+3*float64(config.UIScaling),
		)

		screen.DrawImage(itemTex, itemTexOpts)
	}

	selectedSlotTex := asset_loader.Texture("selected_slot").Texture()
	selectedSlotTexOpts := &ebiten.DrawImageOptions{}
	selectedSlotTexOpts.GeoM.Scale(float64(config.UIScaling), float64(config.UIScaling))
	selectedSlotTexOpts.GeoM.Translate(
		ix+float64(config.UIScaling)+(20*float64(inv.SelectedSlot)*float64(config.UIScaling)),
		iy,
	)
	screen.DrawImage(selectedSlotTex, selectedSlotTexOpts)

	// draw inventory badges on top of everything, so they will be always visible
	inventoryBadgesTex := asset_loader.Texture("inventory_badges").Texture()
	screen.DrawImage(inventoryBadgesTex, inventoryDrawOpts)
}
