package types

import (
	"fmt"
)

type CraftIngredient struct {
	Type   ItemType
	Amount int
}

type Craft struct {
	Name        string
	Description string
	Ingredients []CraftIngredient
	Results     []CraftIngredient
}

// Returns true if the player is able to craft this item
func (craft Craft) AbleToCraft() bool {
	// Check if player has all the necessary ingredients
	for _, item := range craft.Ingredients {
		if !GetInventory().HasItemOfType(item.Type, item.Amount) {
			return false
		}
	}
	// Check if player has space in inventory to store the craft results
	for _, result := range craft.Results {
		if !GetInventory().CanAddItem(ItemSlot{
			Item:     NewItem(result.Type),
			Quantity: uint8(result.Amount),
		}) {
			return false
		}
	}
	return true
}

// Returns true if the item was successfully crafted
func (craft Craft) Craft() bool {
	// First check if there is enough space in inventory to store the craft results
	for _, result := range craft.Results {
		if !GetInventory().CanAddItem(ItemSlot{
			Item:     NewItem(result.Type),
			Quantity: uint8(result.Amount),
		}) {
			return false
		}
	}

	for _, ingredient := range craft.Ingredients {
		removed := GetInventory().RemoveItemByType(ingredient.Type, ingredient.Amount)
		if !removed {
			return false
		}
	}

	for _, result := range craft.Results {
		GetInventory().AddItem(ItemSlot{
			Item:     NewItem(result.Type),
			Quantity: uint8(result.Amount),
		})
	}
	return true
}

// Formats the ingredient list in the following format:
//
// 1. 3x Stick
// 2. 5x Stone
// 3. 10x Sand
func (craft Craft) FormatIngredients() (result string) {
	for i, ingredient := range craft.Ingredients {
		result += fmt.Sprintf("%v. %vx %v\n",
			i+1,
			ingredient.Amount,
			NewItem(ingredient.Type).Name(),
		)
	}
	return
}

// Formats the result list in the following format:
//
// 1. 3x Stick
// 2. 5x Stone
// 3. 10x Sand
func (craft Craft) FormatResults() (result string) {
	for i, ingredient := range craft.Results {
		result += fmt.Sprintf("%v. %vx %v\n",
			i+1,
			ingredient.Amount,
			NewItem(ingredient.Type).Name(),
		)
	}
	return
}
