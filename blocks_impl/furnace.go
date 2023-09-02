package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// After an item has been smelted, wait for 5s until smelting another item
const FurnaceSmeltingCooldown = 300

func init() {
	gob.Register(FurnaceBlockState{})
	types.NewFurnaceBlock = NewFurnaceBlock
}

type FurnaceBlockState struct {
	InputInventory   types.SavedSlot
	OutputInventory  types.SavedSlot
	Energy           float64
	SmeltingCooldown int
}

type FurnaceBlock struct {
	baseBlock
	texturedBlock

	inputInventory  types.ItemSlot
	outputInventory types.ItemSlot

	energy           float64
	smeltingCooldown int // in ticks
}

func (furnace *FurnaceBlock) isSmelting() bool {
	return !furnace.inputInventory.Empty && furnace.energy >= furnace.inputInventory.Item.(types.ISmeltableItem).SmeltingEnergyRequired()
}

func (furnace *FurnaceBlock) updateTexture() {
	if furnace.isSmelting() {
		furnace.tex = assets.Texture("furnace_burning")
	} else {
		furnace.tex = assets.Texture("furnace")
	}
}

func NewFurnaceBlock() types.Block {
	furnace := &FurnaceBlock{
		baseBlock: baseBlock{
			blockType: types.FurnaceBlock,
		},
		texturedBlock: texturedBlock{},

		inputInventory:  types.ItemSlot{Empty: true},
		outputInventory: types.ItemSlot{Empty: true},
	}
	furnace.updateTexture()
	return furnace
}

func (furnace *FurnaceBlock) smeltItem(item types.Item) bool {
	// Accept only smeltable items
	_, isSmeltable := item.(types.ISmeltableItem)
	if !isSmeltable {
		return false
	}

	// If the furnace is empty, add the item straight away
	if furnace.inputInventory.Empty {
		furnace.inputInventory = types.ItemSlot{
			Item:     item,
			Quantity: 1,
		}
		return true
	}

	// Try to add item to the furnace inventory
	return furnace.inputInventory.AddItem(types.ItemSlot{
		Item:     item,
		Quantity: 1,
	})
}

func (furnace *FurnaceBlock) Update(_ types.World) {
	if !furnace.isSmelting() {
		return
	}

	if furnace.smeltingCooldown > 0 {
		furnace.smeltingCooldown--
		return
	}

	if !furnace.inputInventory.Empty {
		furnace.outputInventory.AddItem(types.ItemSlot{
			Item:     furnace.inputInventory.Item.(types.ISmeltableItem).Smelt(),
			Quantity: 1,
		})
		furnace.inputInventory.RemoveItem(1)
		furnace.energy -= furnace.inputInventory.Item.(types.ISmeltableItem).SmeltingEnergyRequired()
		furnace.smeltingCooldown = FurnaceSmeltingCooldown
		furnace.updateTexture()
		furnace.parentChunk.MarkAsModified()
	}
}

func (furnace *FurnaceBlock) Interact() {
	scene_manager.ShowOverlay(&FurnaceScene{furnace})
}

func (furnace *FurnaceBlock) ToolRequiredToBreak() types.ToolFamily {
	return types.ToolFamilyPickaxe
}
func (furnace *FurnaceBlock) ToolStrengthRequired() types.ToolStrength {
	return types.ToolStrengthWood
}
func (furnace *FurnaceBlock) Break() {
	furnace.parentChunk.SetBlock(furnace.x, furnace.y, types.NewGrassBlock())
}

func (furnace *FurnaceBlock) State() interface{} {
	return FurnaceBlockState{
		InputInventory:   furnace.inputInventory.Save(),
		OutputInventory:  furnace.outputInventory.Save(),
		Energy:           furnace.energy,
		SmeltingCooldown: furnace.smeltingCooldown,
	}
}
func (furnace *FurnaceBlock) LoadState(state interface{}) {
	if furnaceState, ok := state.(FurnaceBlockState); ok {
		furnace.inputInventory = furnaceState.InputInventory.Load()
		furnace.outputInventory = furnaceState.OutputInventory.Load()
		furnace.energy = furnaceState.Energy
		furnace.smeltingCooldown = furnaceState.SmeltingCooldown
	}
}

type FurnaceScene struct {
	*FurnaceBlock
}

func (scene *FurnaceScene) Update() {
	// If a player is more than 3 blocks away from the furnace, close the interface
	if scene.FurnaceBlock.Coords().DistanceTo(types.GetCurrentPlayer().Position()) > 3 {
		scene_manager.HideOverlay()
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		if types.GetPlayerInventory().SelectedSlot().Empty {
			break
		}

		itemInHand := types.GetPlayerInventory().ItemInHand()
		if burnableItem, burnable := itemInHand.(types.IBurnableItem); burnable {
			scene.FurnaceBlock.energy += burnableItem.BurningEnergy()
			types.GetPlayerInventory().SelectedSlot().RemoveItem(1)
		} else if scene.FurnaceBlock.smeltItem(itemInHand) {
			types.GetPlayerInventory().SelectedSlot().RemoveItem(1)
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyT):
		if types.GetPlayerInventory().AddItem(scene.FurnaceBlock.outputInventory) {
			scene.FurnaceBlock.outputInventory = types.ItemSlot{Empty: true}
		}
	}
}

func (scene *FurnaceScene) Draw(screen *ebiten.Image) {
	ui.ImmediateDraw(screen, ui.Styled(ui.Padding(1.0,
		ui.VStack().WithSpacing(2.0).WithChildren(
			ui.Background(colors.C("blue"), ui.PaddingXY(1.0, 0.3,
				ui.CustomLabel("Furnace", colors.C("white"), 1.5),
			)),

			ui.Tooltip(ui.VStack().WithSpacing(1.5).WithChildren(
				ui.VStack(
					ui.Label("Fuel"),
					ui.LabelF("%v", scene.FurnaceBlock.energy),
				),
				ui.VStack(
					ui.Label("Items to smelt"),
					ui.ItemSlot(&scene.FurnaceBlock.inputInventory),
				),
				ui.VStack(
					ui.Label("Smelted items"),
					ui.ItemSlot(&scene.FurnaceBlock.outputInventory),
				),
			)),
		),
	)).WithTextColor(colors.C("white")))
}

func (scene *FurnaceScene) Destroy() {

}
