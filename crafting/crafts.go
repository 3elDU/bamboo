package crafting

import "github.com/3elDU/bamboo/types"

// A list of all available crafts
var Crafts []types.Craft = []types.Craft{
	{
		Name: "1",
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.StickItem,
				Amount: 1,
			},
		},
		Result: []types.CraftIngredient{
			{
				Type:   types.FlintItem,
				Amount: 1,
			},
		},
	},

	{
		Name: "2",
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.StickItem,
				Amount: 1,
			},
		},
		Result: []types.CraftIngredient{
			{
				Type:   types.FlintItem,
				Amount: 2,
			},
		},
	},
}

// Returns true if the player has necessary ingredients to craft this item
func AbleToCraft(craft types.Craft) bool {
	for _, item := range craft.Ingredients {
		if !types.GetInventory().HasItemOfType(item.Type, item.Amount) {
			return false
		}
	}
	return true
}

func Craft(craft types.Craft) bool {
	for _, ingredient := range craft.Ingredients {
		removed := types.GetInventory().RemoveItemByType(ingredient.Type, ingredient.Amount)
		if !removed {
			return false
		}
	}

	for _, result := range craft.Result {
		types.GetInventory().AddItem(types.ItemSlot{
			Item:     types.NewItem(result.Type),
			Quantity: uint8(result.Amount),
		})
	}
	return true
}
