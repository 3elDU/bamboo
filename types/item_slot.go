package types

import "github.com/3elDU/bamboo/config"

// Holds multiple items of the same type
type ItemSlot struct {
	Item     Item
	Quantity uint8
	Empty    bool
}

// Returns true if item has been successfully added
// False if there is no space, or item is of different type
func (slot *ItemSlot) AddItem(other ItemSlot) bool {
	// If a slot we're trying to add is empty, do not do anything
	if other.Empty {
		return true
	}

	if slot.Empty {
		slot.Item = other.Item
		slot.Quantity = other.Quantity
		slot.Empty = false
		return true
	}

	if other.Item.Hash() != slot.Item.Hash() {
		return false
	}

	if slot.Quantity+other.Quantity > config.SlotSize {
		return false
	}

	slot.Quantity += other.Quantity
	return true
}

func (slot *ItemSlot) RemoveItem(count uint8) {
	if slot.Quantity <= count {
		slot.Quantity = 0
	} else {
		slot.Quantity -= count
	}
}
