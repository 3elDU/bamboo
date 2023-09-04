package inventory

import (
	"fmt"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

const Size = 5

type Inventory struct {
	Slots        [Size]*types.ItemSlot
	selectedSlot int
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	for i := range inv.Slots {
		inv.Slots[i] = new(types.ItemSlot)
		inv.Slots[i].Empty = true
	}

	types.SetPlayerInventory(inv)
	return inv
}

func (inv *Inventory) Length() int {
	return Size
}

func (inv *Inventory) At(i int) *types.ItemSlot {
	return inv.Slots[i]
}

func (inv *Inventory) HasItemOfType(itemType types.ItemType, amount int) bool {
	for _, slot := range inv.Slots {
		if slot.Empty {
			continue
		}

		if slot.Item.Type() == itemType && slot.Quantity >= uint8(amount) {
			return true
		}
	}
	return false
}

func (inv *Inventory) RemoveItem(item types.ItemSlot) bool {
	for i, slot := range inv.Slots {
		if slot.Empty {
			continue
		}

		if slot.Item.Stackable() && item.Item.Type() == slot.Item.Type() && slot.Quantity >= item.Quantity {
			slot.RemoveItem(item.Quantity)
			if slot.Quantity == 0 {
				inv.Slots[i] = new(types.ItemSlot)
				inv.Slots[i].Empty = true
			}
			return true
		}
	}
	return false
}

func (inv *Inventory) RemoveItemByType(itemType types.ItemType, amount int) bool {
	for i, slot := range inv.Slots {
		if slot.Empty {
			continue
		}

		if slot.Item.Type() == itemType && slot.Quantity >= uint8(amount) {
			slot.RemoveItem(uint8(amount))
			if slot.Quantity == 0 {
				inv.Slots[i] = new(types.ItemSlot)
				inv.Slots[i].Empty = true
			}
			return true
		}
	}
	return false
}

func (inv *Inventory) CanAddItems(items ...types.ItemSlot) bool {
	for _, item := range items {
		if !inv.CanAddItem(item) {
			return false
		}
	}
	return true
}
func (inv *Inventory) CanAddItem(item types.ItemSlot) bool {
	// Skip the check if this is an empty slot
	if item.Empty {
		return true
	}

	for _, slot := range inv.Slots {
		if slot.Empty {
			return true
		}
		if slot.Item.Type() == item.Item.Type() && slot.Quantity+item.Quantity <= config.SlotSize {
			return true
		}
	}
	return false
}

func (inv *Inventory) AddItem(item types.ItemSlot) bool {
	// Skip the check if this is an empty slot
	if item.Empty {
		return true
	}

	// Try to find a slot that already holds item with the same type
	for _, slot := range inv.Slots {
		if !slot.Empty && slot.Item.Type() == item.Item.Type() && slot.Item.Stackable() {
			if slot.AddItem(item) {
				return true
			}
		}
	}

	for _, slot := range inv.Slots {
		if slot.AddItem(item) {
			return true
		}
	}

	return false
}
func (inv *Inventory) AddItems(items ...types.ItemSlot) bool {
	// First check if all the items can be added to the inventory
	if !inv.CanAddItems(items...) {
		return false
	}

	for _, item := range items {
		inv.AddItem(item)
	}
	return true
}

func (inv *Inventory) SelectSlot(slot int) {
	if slot >= Size {
		slot = 0
	} else if slot < 0 {
		slot = Size - 1
	}

	inv.selectedSlot = slot
}

func (inv *Inventory) ItemInHand() types.Item {
	return inv.Slots[inv.selectedSlot].Item
}
func (inv *Inventory) SelectedSlot() *types.ItemSlot {
	return inv.Slots[inv.selectedSlot]
}
func (inv *Inventory) SelectedSlotIndex() int {
	return inv.selectedSlot
}

// Returns a position of inventory slot on the screen
func (inv *Inventory) SlotToScreenCoords(screen *ebiten.Image, slot int) types.Vec2f {
	inventoryTexture := assets.Texture("inventory")
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
	inventoryTexture := assets.Texture("inventory")
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

		// Render label with item amount only if there is more than 1 of that item
		if slot.Quantity > 1 {
			font.RenderFont(screen, fmt.Sprintf("%v", slot.Quantity), itemTexPos.X, itemTexPos.Y, colors.C("white"))
		}
	}

	// Draw an outline around the selected slot
	selectedSlotTex := assets.Texture("selected_slot").Texture()
	selectedSlotTexOpts := &ebiten.DrawImageOptions{}
	selectedSlotTexOpts.GeoM.Scale(config.UIScaling, config.UIScaling)
	selectedSlotTexOpts.GeoM.Translate(
		ix+config.UIScaling+(20*float64(inv.selectedSlot)*config.UIScaling),
		iy,
	)
	screen.DrawImage(selectedSlotTex, selectedSlotTexOpts)

	// Draw inventory badges on top of everything, so they will be always visible
	inventoryBadgesTex := assets.Texture("inventory_badges").Texture()
	screen.DrawImage(inventoryBadgesTex, inventoryDrawOpts)

	// Check if cursor hovers over one of the items in inventory, and render item's tooltip
	for i := 0; i < inv.Length(); i++ {
		slot := inv.At(i)
		if slot.Empty {
			continue
		}

		if inv.MouseOverSlot(screen, i) {
			item := slot.Item

			var tooltipText string
			if item.Description() == "" {
				tooltipText = item.Name()
			} else {
				tooltipText = fmt.Sprintf("%v\n------\n%v", item.Name(), item.Description())
			}

			cx, cy := ebiten.CursorPosition()
			ui.DrawTextTooltip(screen, cx, cy, ui.TopRight, tooltipText)
		}
	}
}
