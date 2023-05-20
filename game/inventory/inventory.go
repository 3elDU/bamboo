package inventory

import (
	"fmt"
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

const Size = 5

type Inventory struct {
	Slots        [Size]*types.ItemSlot
	SelectedSlot int
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	for i := range inv.Slots {
		inv.Slots[i] = new(types.ItemSlot)
		inv.Slots[i].Empty = true
	}

	types.SetInventory(inv)
	return inv
}

func (inv *Inventory) Length() int {
	return Size
}

func (inv *Inventory) At(i int) types.ItemSlot {
	return *inv.Slots[i]
}

func (inv *Inventory) RemoveItem(item types.ItemSlot) bool {
	for i, slot := range inv.Slots {
		if slot.Empty || slot.Item.Hash() != item.Item.Hash() {
			continue
		}

		slot.RemoveItem(item.Quantity)
		if slot.Quantity == 0 {
			inv.Slots[i] = new(types.ItemSlot)
			inv.Slots[i].Empty = true
		}
		return true
	}
	return false
}

func (inv *Inventory) AddItem(item types.ItemSlot) bool {
	for _, slot := range inv.Slots {
		if slot.AddItem(item) {
			return true
		}
	}
	return false
}

func (inv *Inventory) SelectSlot(slot int) {
	if slot >= Size {
		slot = 0
	} else if slot < 0 {
		slot = Size - 1
	}

	inv.SelectedSlot = slot
}

func (inv *Inventory) ItemInHand() types.Item {
	return inv.Slots[inv.SelectedSlot].Item
}

// Returns a position of inventory slot on the screen
func (inv *Inventory) SlotToScreenCoords(screen *ebiten.Image, slot int) types.Vec2f {
	inventoryTexture := asset_loader.Texture("inventory")
	w, h := inventoryTexture.ScaledSize()
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	return types.Vec2f{
		X: (float64(sw)/2 - float64(w)/2) + 4*float64(config.UIScaling) + (20 * float64(slot) * config.UIScaling),
		Y: (float64(sh) - h) + 3*float64(config.UIScaling),
	}
}

func (inv *Inventory) MouseOverSlot(screen *ebiten.Image, slot int) bool {
	itemTexPos := inv.SlotToScreenCoords(screen, slot)
	cx, cy := ebiten.CursorPosition()
	return float64(cx) > itemTexPos.X && float64(cy) > itemTexPos.Y && float64(cx) < itemTexPos.X+16*config.UIScaling && float64(cy) < itemTexPos.Y+16*config.UIScaling
}

func (inv *Inventory) Render(screen *ebiten.Image) {
	inventoryTexture := asset_loader.Texture("inventory")
	inventoryDrawOpts := &ebiten.DrawImageOptions{}

	w, h := inventoryTexture.ScaledSize()
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()

	// position of inventory texture on the screen
	ix := float64(sw)/2 - float64(w)/2 // horizontally centered
	iy := float64(sh) - h              // bottom of the screen

	inventoryDrawOpts.GeoM.Scale(config.UIScaling, config.UIScaling)
	inventoryDrawOpts.GeoM.Translate(ix, iy)

	screen.DrawImage(inventoryTexture.Texture(), inventoryDrawOpts)

	for i, slot := range inv.Slots {
		if slot.Empty {
			continue
		}

		itemTex := slot.Item.Texture()
		itemTexOpts := &ebiten.DrawImageOptions{}

		itemTexPos := inv.SlotToScreenCoords(screen, i)

		itemTexOpts.GeoM.Scale(config.UIScaling, config.UIScaling)
		itemTexOpts.GeoM.Translate(itemTexPos.X, itemTexPos.Y)

		screen.DrawImage(itemTex, itemTexOpts)

		font.RenderFont(screen, fmt.Sprintf("%v", slot.Quantity), itemTexPos.X, itemTexPos.Y, colors.Black)
	}

	selectedSlotTex := asset_loader.Texture("selected_slot").Texture()
	selectedSlotTexOpts := &ebiten.DrawImageOptions{}
	selectedSlotTexOpts.GeoM.Scale(config.UIScaling, config.UIScaling)
	selectedSlotTexOpts.GeoM.Translate(
		ix+config.UIScaling+(20*float64(inv.SelectedSlot)*config.UIScaling),
		iy,
	)
	screen.DrawImage(selectedSlotTex, selectedSlotTexOpts)

	// draw inventory badges on top of everything, so they will be always visible
	inventoryBadgesTex := asset_loader.Texture("inventory_badges").Texture()
	screen.DrawImage(inventoryBadgesTex, inventoryDrawOpts)
}
