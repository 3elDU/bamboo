package types

import "github.com/3elDU/bamboo/config"

// Holds multiple items of the same type
type ItemSlot struct {
	Item     Item
	Quantity uint8
	Empty    bool
}

func NewItemSlot(item Item, quantity uint8) ItemSlot {
	var empty bool = false
	if quantity == 0 {
		empty = true
	}
	return ItemSlot{
		Item:     item,
		Quantity: quantity,
		Empty:    empty,
	}
}

func (slot *ItemSlot) CanAddItem(other ItemSlot) bool {
	if slot.Empty || other.Empty {
		return true
	}

	if slot.Item.Type() != other.Item.Type() {
		return false
	}

	return slot.Quantity+other.Quantity <= config.SlotSize
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

	if other.Item.Type() != slot.Item.Type() {
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

	if slot.Quantity == 0 {
		slot.Empty = true
	}
}

func (slot *ItemSlot) Save() SavedSlot {
	var itemType ItemType
	var itemState interface{}
	if !slot.Empty {
		itemType = slot.Item.Type()
		itemState = slot.Item.State()
	}

	return SavedSlot{
		Empty:    slot.Empty,
		Quantity: slot.Quantity,
		ItemType: itemType,
		State:    itemState,
	}
}
