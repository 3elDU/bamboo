package game

import (
	"image/color"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/crafting"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/exp/slices"
)

type craftingMenu struct {
	availableCrafts []types.Craft
	selectedCraft   int

	craftList                *ui.AvailableCraftsList
	selectedCraftDescription *ui.CraftDescription
}

func newCraftingMenu() *craftingMenu {
	menu := &craftingMenu{selectedCraft: 0}
	menu.craftList = ui.NewAvailableCraftsList(menu.availableCrafts, menu.selectedCraft)
	menu.selectedCraftDescription = ui.NewCraftDescription()
	return menu
}

func (menu *craftingMenu) SelectedCraft() types.Craft {
	if len(menu.availableCrafts) <= menu.selectedCraft {
		return types.Craft{}
	}
	return menu.availableCrafts[menu.selectedCraft]
}

func (menu *craftingMenu) setSelectedCraft(selectedCraft int) {
	if selectedCraft < 0 {
		selectedCraft = len(menu.availableCrafts) - 1
	}
	if selectedCraft > len(menu.availableCrafts)-1 {
		selectedCraft = 0
	}
	menu.selectedCraft = selectedCraft
	menu.craftList.SetSelectedCraft(selectedCraft)
	menu.selectedCraftDescription.SetCraft(menu.SelectedCraft())
}

// Updates the list of recipes that can be crafted, and updated the UI accordingly.
func (menu *craftingMenu) UpdateAvailableRecipes() {
	menu.availableCrafts = []types.Craft{}
	for _, recipe := range crafting.Crafts {
		if recipe.AbleToCraft() {
			menu.availableCrafts = append(menu.availableCrafts, recipe)
		}
	}

	// Sort the list of available crafts alphabetically
	slices.SortStableFunc(menu.availableCrafts, func(a types.Craft, b types.Craft) bool {
		for _, charA := range a.Name {
			for _, charB := range b.Name {
				if charA < charB {
					return true
				}
			}
		}
		return false
	})

	menu.craftList.SetAvailableCrafts(menu.availableCrafts)
	if len(menu.availableCrafts) > 0 {
		menu.setSelectedCraft(menu.selectedCraft)
	}
}

func (menu *craftingMenu) Update() {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		menu.SelectedCraft().Craft()
		menu.UpdateAvailableRecipes()

		// If crafting with a shift key, do not close the crafting menu
		if !ebiten.IsKeyPressed(ebiten.KeyShift) {
			scene_manager.HideOverlay()
		}
	}

	// Return early if there are no available crafts.
	// Handling up/down arrow keys makes no sense, when there is no items to choose from
	if len(menu.availableCrafts) == 0 {
		scene_manager.HideOverlay()
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		menu.setSelectedCraft(menu.selectedCraft + 1)
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		menu.setSelectedCraft(menu.selectedCraft - 1)
	}
}

func (menu *craftingMenu) Draw(screen *ebiten.Image) {
	var view ui.Component

	if len(menu.availableCrafts) > 0 {
		view = ui.Screen(
			ui.Padding(1.0, ui.VStack().WithSpacing(2.0).WithChildren(
				ui.Background(colors.C("blue"),
					ui.PaddingXY(1.0, 0.3, ui.CustomLabel("Crafting", colors.C("white"), 1.5)),
				),
				ui.HStack().WithSpacing(2.0).WithChildren(
					ui.VStack(
						ui.ColoredLabel("List", color.White),
						ui.Tooltip(menu.craftList),
					),
					ui.VStack(
						ui.ColoredLabel("Description", color.White),
						ui.Tooltip(menu.selectedCraftDescription),
					),
				),
			)))
	} else {
		// If there are no available crafts, show a placeholder text
		view = ui.Screen(ui.Padding(1.0,
			ui.Tooltip(ui.Padding(1.0,
				ui.ColoredLabel("No crafts available!", color.White),
			)),
		))
	}

	view.Draw(screen, 0, 0)
}

func (menu *craftingMenu) Destroy() {

}
