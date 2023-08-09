package crafting

import "github.com/3elDU/bamboo/types"

// A list of all available crafts
var Crafts []types.Craft = []types.Craft{
	{
		Name:        "Watering can",
		Description: "Water your crops!",
		Conditions:  []types.CraftCondition{PlayerMustBeNearCampfire},
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.ClayItem,
				Amount: 2,
			},
		},
		Result: types.CraftIngredient{
			Type:   types.WateringCanItem,
			Amount: 1,
		},
	},
	{
		Name:       "Clay shovel",
		Conditions: []types.CraftCondition{PlayerMustBeNearCampfire},
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.ClayItem,
				Amount: 3,
			},
		},
		Result: types.CraftIngredient{
			Type:   types.ClayShovelItem,
			Amount: 1,
		},
	},
}
