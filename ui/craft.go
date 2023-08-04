package ui

import (
	"image/color"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/types"
)

type AvailableCraftsList struct {
	*StackComponent

	availableCrafts []types.Craft
	selectedCraft   int
}

func NewAvailableCraftsList(availableCrafts []types.Craft, selectedCraft int) *AvailableCraftsList {
	craftList := &AvailableCraftsList{
		StackComponent: nil,

		availableCrafts: []types.Craft{},
		selectedCraft:   -1,
	}
	craftList.SetAvailableCrafts(availableCrafts)
	craftList.SetSelectedCraft(selectedCraft)
	return craftList
}

// Update the list of available crafts.
// Also updates the UI
func (craftList *AvailableCraftsList) SetAvailableCrafts(availableCrafts []types.Craft) {
	craftList.availableCrafts = availableCrafts

	stack := VStack().AlignChildren(AlignCenter)
	for _, craft := range availableCrafts {
		stack.AddChild(
			PaddingX(1.0, ColoredLabel(craft.Name, color.White)),
		)
	}
	craftList.StackComponent = stack
}

func (craftList *AvailableCraftsList) SetSelectedCraft(selectedCraft int) {
	// reassemble the stack
	craftList.SetAvailableCrafts(craftList.availableCrafts)

	for i, child := range craftList.Children() {
		if i == selectedCraft {
			craftList.ReplaceChild(child,
				Background(colors.C("blue"),
					PaddingX(1.0,
						ColoredLabel(craftList.availableCrafts[i].Name, color.White),
					),
				),
			)
		}
	}
}

type CraftDescription struct {
	*StyledComponent

	craft types.Craft
}

func NewCraftDescription() *CraftDescription {
	craftDescription := &CraftDescription{
		StyledComponent: Styled(HStack()),
		craft:           types.Craft{},
	}
	return craftDescription
}

func (craftDescription *CraftDescription) SetCraft(craft types.Craft) {
	craftDescription.craft = craft

	ingredientsStack := VStack().WithSpacing(0.2).WithChildren(
		Label("Ingredients:"),
	)
	for _, ingredient := range craft.Ingredients {
		ingredientsStack.AddChild(HStack(
			Tooltip(Image(types.NewItem(ingredient.Type).Texture())).
				WithNeutralColor(),
			LabelF("%vx", ingredient.Amount),
			Label(types.NewItem(ingredient.Type).Name()),
		).WithSpacing(1).AlignChildren(AlignCenter))
	}

	craftDescription.StyledComponent = Styled(
		Padding(0.3, VStack().WithSpacing(1.0).WithChildren(
			// Header
			HStack(
				Tooltip(Image(types.NewItem(craft.Result.Type).Texture())).
					WithNeutralColor(),
				Label(types.NewItem(craft.Result.Type).Name()),
			).WithSpacing(1).AlignChildren(AlignCenter),

			// Description
			Label(craft.Description),

			// Ingredients
			ingredientsStack,
		)),
	).WithTextColor(color.White)
}
