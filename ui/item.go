package ui

import (
	"fmt"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type ItemSlotComponent struct {
	*TooltipComponent

	itemSlot *types.ItemSlot
}

func ItemSlot(slot *types.ItemSlot) *ItemSlotComponent {
	var itemCountLabel = ""
	// Show item count only if there's more than 1 item
	if slot.Quantity > 1 {
		itemCountLabel = fmt.Sprint(slot.Quantity)
	}

	var overlay *OverlayComponent
	if slot.Empty {
		overlay = Overlay(
			Image(assets.Texture("empty").Texture()),
		)
	} else {
		overlay = Overlay(
			Image(slot.Item.Texture()),
			Label(itemCountLabel),
		)
	}

	tooltip := Tooltip(overlay).WithNeutralColor()

	return &ItemSlotComponent{
		TooltipComponent: tooltip,
		itemSlot:         slot,
	}
}

func (slot *ItemSlotComponent) Update() error {
	return slot.TooltipComponent.Update()
}
func (slot *ItemSlotComponent) Draw(screen *ebiten.Image, x, y float64) error {
	if err := slot.TooltipComponent.Draw(screen, x, y); err != nil {
		return err
	}

	w, h := slot.TooltipComponent.ComputedSize()
	cx, cy := ebiten.CursorPosition()

	if cx >= int(x) && cy >= int(y) && cx <= int(x+w) && cy <= int(y+h) {
		if slot.itemSlot.Empty {
			DrawTextTooltip(screen, cx, cy, TopRight, "(empty)")
		} else {
			DrawTextTooltip(screen, cx, cy, TopRight, fmt.Sprintf(
				"%vx %v",
				slot.itemSlot.Quantity, slot.itemSlot.Item.Name(),
			))
		}
	}
	return nil
}
