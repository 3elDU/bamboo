package crafting

import "github.com/3elDU/bamboo/types"

// A list of all available crafts
var Crafts []types.Craft = []types.Craft{
	{
		Name:        "Debug1",
		Description: "Make a stick from a flint!\nHow convenient!",
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.FlintItem,
				Amount: 1,
			},
		},
		Results: []types.CraftIngredient{
			{
				Type:   types.StickItem,
				Amount: 1,
			},
		},
	},

	{
		Name:        "Debug2",
		Description: "Make a flint from a stick!\nHow convenient!",
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.StickItem,
				Amount: 1,
			},
		},
		Results: []types.CraftIngredient{
			{
				Type:   types.FlintItem,
				Amount: 1,
			},
		},
	},

	{
		Name:        "Test3",
		Description: "lorem ipsum dolor sit amet...",
		Ingredients: []types.CraftIngredient{
			{
				Type:   types.FlintItem,
				Amount: 30,
			},
		},
		Results: []types.CraftIngredient{
			{
				Type:   types.StickItem,
				Amount: 30,
			},
		},
	},
}
