package types

var currentInventory Inventory

// Sets the current instance of Inventory
func SetPlayerInventory(inventory Inventory) {
	currentInventory = inventory
}

// Returns the current instance of Inventory,
// so that inventory can be accessed from everywhere in the code
func GetPlayerInventory() Inventory {
	return currentInventory
}

type Inventory interface {
	// Returns the number of total slots
	Length() int
	At(i int) *ItemSlot

	// Checks if there is enough space in the inventory to add this item
	CanAddItem(item ItemSlot) bool
	CanAddItems(items ...ItemSlot) bool

	// AddItem returns false if there is no space
	AddItem(item ItemSlot) bool
	AddItems(items ...ItemSlot) bool

	RemoveItem(item ItemSlot) bool
	RemoveItemByType(itemType ItemType, amount int) bool

	SelectSlot(i int)
	ItemInHand() Item
	SelectedSlot() *ItemSlot
	SelectedSlotIndex() int

	HasItemOfType(itemType ItemType, amount int) bool
}

type SavedSlot struct {
	Empty    bool
	Quantity uint8
	ItemType ItemType
	State    interface{}
}

func (savedSlot *SavedSlot) Load() ItemSlot {
	if savedSlot.Empty {
		return ItemSlot{Empty: true}
	}

	slot := ItemSlot{}
	slot.Item = NewItem(savedSlot.ItemType)
	slot.Quantity = savedSlot.Quantity
	slot.Item.LoadState(savedSlot.State)
	return slot
}
