package inventory

import (
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

const InventorySize = 5

type Inventory struct {
	Slots        [InventorySize]*types.ItemSlot
	SelectedSlot uint8
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

func (inv *Inventory) SelectSlot(slot uint8) {
	if slot > InventorySize {
		inv.SelectedSlot = InventorySize - 1
	}

	inv.SelectedSlot = slot
}

func (inv *Inventory) ItemInHand() types.Item {
	return inv.Slots[inv.SelectedSlot].Item
}

func (inv *Inventory) Render(screen *ebiten.Image) {
	tex := asset_loader.Texture("inventory").Texture()
	badge_tex := asset_loader.Texture("inventory_badges").Texture()
	w, h := tex.Size()
	opts := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()

	// position of inventory texture on the screen
	ix := float64(sw)/2 - float64(w)/2
	iy := float64(sh) - float64(h)

	opts.GeoM.Translate(ix, iy)

	screen.DrawImage(tex, opts)

	for i, slot := range inv.Slots {
		if slot.Empty {
			continue
		}

		itemTex := slot.Item.Texture()
		itemTexOpts := &ebiten.DrawImageOptions{}

		itemTexOpts.GeoM.Scale(2, 2)
		itemTexOpts.GeoM.Translate(ix+6+36*float64(i), iy+4)

		screen.DrawImage(itemTex, itemTexOpts)
	}

	selectedSlotTex := asset_loader.Texture("selected_slot").Texture()
	selectedSlotTexOpts := &ebiten.DrawImageOptions{}
	selectedSlotTexOpts.GeoM.Translate(ix+2+36*float64(inv.SelectedSlot), iy)
	screen.DrawImage(selectedSlotTex, selectedSlotTexOpts)

	screen.DrawImage(badge_tex, opts)
}
