package types

type CraftIngredient struct {
	Type   ItemType
	Amount int
}

type Craft struct {
	Name        string
	Description string
	Ingredients []CraftIngredient
	Result      []CraftIngredient
}
