package types

type CraftIngredient struct {
	Type   ItemType
	Amount int
}

// CraftCondition is a function that returns true, if a player is able to craft this
// For example, to craft some item, we can check if the player is standing near a specific block.
type CraftCondition func() bool

type Craft struct {
	Name        string
	Description string
	Conditions  []CraftCondition
	Ingredients []CraftIngredient
	Result      CraftIngredient
}

// Returns true if the player is able to craft this item
func (craft Craft) AbleToCraft() bool {
	// Check if all craft conditions all met
	for _, condition := range craft.Conditions {
		if !condition() {
			return false
		}
	}

	// Check if player has all the necessary ingredients
	for _, item := range craft.Ingredients {
		if !GetInventory().HasItemOfType(item.Type, item.Amount) {
			return false
		}
	}
	// Check if player has space in inventory to store the craft result
	if !GetInventory().CanAddItem(ItemSlot{
		Item:     NewItem(craft.Result.Type),
		Quantity: uint8(craft.Result.Amount),
	}) {
		return false
	}

	return true
}

// Returns true if the item was successfully crafted
func (craft Craft) Craft() bool {
	if !craft.AbleToCraft() {
		return false
	}

	for _, ingredient := range craft.Ingredients {
		removed := GetInventory().RemoveItemByType(ingredient.Type, ingredient.Amount)
		if !removed {
			return false
		}
	}

	GetInventory().AddItem(ItemSlot{
		Item:     NewItem(craft.Result.Type),
		Quantity: uint8(craft.Result.Amount),
	})

	return true
}
