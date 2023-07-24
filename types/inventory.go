package types

var currentInventory Inventory

// Sets the current instance of Inventory
func SetInventory(inventory Inventory) {
	currentInventory = inventory
}

// Returns the current instance of Inventory,
// so that inventory can be accessed from everywhere in the code
func GetInventory() Inventory {
	return currentInventory
}

type Inventory interface {
	// Returns the number of total slots
	Length() int
	At(i int) ItemSlot

	// Checks if there is enough space in the inventory to add this item
	CanAddItem(item ItemSlot) bool
	// AddItem returns false if there is no space
	AddItem(item ItemSlot) bool
	RemoveItem(item ItemSlot) bool
	RemoveItemByType(itemType ItemType, amount int) bool

	SelectSlot(i int)
	ItemInHand() Item

	HasItemOfType(itemType ItemType, amount int) bool
}
