package game

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/3elDU/bamboo/crafting"
	"github.com/3elDU/bamboo/game/inventory"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type craftingMenu struct {
	selectedRecipe int
	inventory      *inventory.Inventory

	stack *ui.StackComponent
}

func newCraftingMenu(inventory *inventory.Inventory) *craftingMenu {
	stack := ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 1.0})
	for _, craft := range crafting.Crafts {
		if crafting.AbleToCraft(craft) {
			stack.AddChild(ui.Label(ui.LabelOptions{Color: color.White, Scaling: 1.0}, craft.Name))
		}
	}

	return &craftingMenu{
		inventory: inventory,
		stack:     stack,
	}
}

func (menu *craftingMenu) Update() bool {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		crafting.Craft(crafting.Crafts[menu.selectedRecipe])
		return true

	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		menu.selectedRecipe += 1
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		menu.selectedRecipe -= 1
	}
	return false
}

func (menu *craftingMenu) Draw(screen *ebiten.Image) {
	for i, child := range menu.stack.Children() {
		textView := child.(ui.TextView)
		text := textView.Text()

		if i == menu.selectedRecipe && !strings.Contains(text, "-> ") {
			text = "-> " + text
		} else if i != menu.selectedRecipe {
			text = strings.ReplaceAll(text, "-> ", "")
		}
		textView.SetText(text)
	}

	craftInfo := ""
	craftInfo += fmt.Sprintf("%v\n%v\n", crafting.Crafts[menu.selectedRecipe].Name, crafting.Crafts[menu.selectedRecipe].Description)
	craftInfo += fmt.Sprintf("------ Ingredients:\n")
	for i, ingredient := range crafting.Crafts[menu.selectedRecipe].Ingredients {
		craftInfo += fmt.Sprintf("%v. %vx %v\n",
			i+1,
			ingredient.Amount,
			types.NewItem(ingredient.Type).Name(),
		)
	}
	craftInfo += fmt.Sprintf("------ Result\n")
	for i, result := range crafting.Crafts[menu.selectedRecipe].Result {
		craftInfo += fmt.Sprintf("%v. %vx %v\n",
			i+1,
			result.Amount,
			types.NewItem(result.Type).Name(),
		)
	}

	view := ui.Screen(
		ui.Padding(1.0, ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 5.0},
			ui.Tooltip(menu.stack),
			ui.Tooltip(ui.Label(ui.LabelOptions{Color: color.White, Scaling: 1.0}, craftInfo)),
		)),
	)
	view.Draw(screen, 0, 0)
}
