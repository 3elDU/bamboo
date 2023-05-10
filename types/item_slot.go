package types

// Holds multiple items of the same type
type ItemSlot struct {
	Item     Item
	Quantity uint8
	Empty    bool
}

// Returns true if item has been successfully added
// False if there is no space, or item is of different type
func (slot *ItemSlot) AddItem(item Item) bool {
	if slot.Empty {
		slot.Item = item
		slot.Quantity = 1
		slot.Empty = false
		return true
	}

	if item.Hash() != slot.Item.Hash() {
		return false
	}

	if slot.Quantity > 50 {
		return false
	}

	slot.Quantity++
	return true
}

func (slot *ItemSlot) RemoveItem(count uint8) {
	if slot.Quantity <= count {
		slot.Quantity = 0
	} else {
		slot.Quantity -= count
	}
}
