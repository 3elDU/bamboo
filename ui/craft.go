package ui

import (
	"image/color"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/types"
	"github.com/MakeNowJust/heredoc"
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
	*LabelComponent

	craft types.Craft
}

func NewCraftDescription(craft types.Craft) *CraftDescription {
	craftDescription := &CraftDescription{
		LabelComponent: ColoredLabel("", color.White),
		craft:          types.Craft{},
	}
	craftDescription.SetCraft(craft)
	return craftDescription
}

func (craftDescription *CraftDescription) SetCraft(craft types.Craft) {
	craftDescription.craft = craft

	// Update label
	craftDescription.SetText(heredoc.Docf(`
		%v
		%v
		-- Ingredients
		%v
		-- Results
		%v
		`, craft.Name, craft.Description, craft.FormatIngredients(), craft.FormatResults(),
	))
}
