package game

import (
	"fmt"
	"log"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/event"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/game/inventory"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/ui"
	"github.com/3elDU/bamboo/world"
	"github.com/3elDU/bamboo/world_type"
	"github.com/MakeNowJust/heredoc"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	paused    bool
	pauseMenu *pauseMenu

	world       *world.World
	player      *player.Player
	playerStack *player.Stack
	inventory   *inventory.Inventory

	craftingMenu *craftingMenu
	compass      *ui.CompassComponent

	// debug switches
	superSpeed bool
}

func newGame(gameWorld *world.World, playerStack *player.Stack, inventory *inventory.Inventory) *Game {
	game := &Game{
		pauseMenu:    newPauseMenu(),
		craftingMenu: newCraftingMenu(),

		world:       gameWorld,
		playerStack: playerStack,
		inventory:   inventory,

		compass: ui.NewCompassComponent(),
	}
	game.player = playerStack.Top()

	return game
}

// Creates a game scene with a new world
func NewGameScene(metadata types.Save) *Game {
	w := world.NewWorld(metadata)
	stack := player.NewPlayerStack()
	stack.Push(player.NewPlayer(w))
	game := newGame(
		w,
		stack,
		inventory.NewInventory(),
	)

	// perform a save immediately after the scene creation
	game.Save()

	return game
}

// Creates a game scene from existing world
func LoadGameScene(metadata types.Save) *Game {
	// load the player stack first, to determine which world to load
	loadedPlayer := player.LoadPlayerStack(metadata.BaseUUID)
	loadedWorld := world.Load(metadata.BaseUUID, loadedPlayer.Top().SelectedWorld.UUID)
	loadedInventory := inventory.LoadInventory(metadata.BaseUUID)
	return newGame(loadedWorld, loadedPlayer, loadedInventory)
}

func (game *Game) Save() {
	game.world.Save()
	game.playerStack.Save(game.world.Metadata())
	game.inventory.Save(game.world.Metadata())
}

func (game *Game) processInput() {
	game.player.UpdateInput(player.MovementVector{
		Left:  ebiten.IsKeyPressed(ebiten.KeyA),
		Right: ebiten.IsKeyPressed(ebiten.KeyD),
		Up:    ebiten.IsKeyPressed(ebiten.KeyW),
		Down:  ebiten.IsKeyPressed(ebiten.KeyS),
	})

	lookingAt := game.player.LookingAt()
	// Check for key presses
	switch {
	// Escape key
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		// If an overlay is being shown, hide it
		if scene_manager.DisplayingOverlay() {
			scene_manager.HideOverlay()
			game.paused = false
		} else {
			// Otherwise, open the pause menu
			scene_manager.ShowOverlay(newPauseMenu())
			game.paused = true
		}

	// Open crafting menu
	case inpututil.IsKeyJustPressed(ebiten.KeyC):
		log.Println("Entering crafting menu")
		game.craftingMenu.UpdateAvailableRecipes()
		scene_manager.ShowOverlay(game.craftingMenu)

	// Break the block
	case inpututil.IsKeyJustPressed(ebiten.KeyR):
		block, breakable := game.world.BlockAt(lookingAt.X, lookingAt.Y).(types.BreakableBlock)
		if !breakable {
			break
		}

		// Check if the tool can break the block
		tool, isTool := game.inventory.ItemInHand().(types.IToolItem)

		// Check if block can be broken with the bare hand
		if block.ToolStrengthRequired() == types.ToolStrengthBareHand {
			block.Break()
			break
		}

		if isTool && tool.ToolFamily() == block.ToolRequiredToBreak() && tool.ToolStrength() >= block.ToolStrengthRequired() {
			block.Break()
		}

	// Use the item in hand / Interact with the block
	case inpututil.IsKeyJustPressed(ebiten.KeyF):
		// If the block is interactive, interact with it
		if block, ok := game.world.BlockAt(lookingAt.X, lookingAt.Y).(types.InteractiveBlock); ok {
			block.Interact()
		} else if tool, itemIsTool := game.inventory.ItemInHand().(types.IToolItem); itemIsTool {
			// Otherwise, use the item in hand
			tool.UseTool(lookingAt)
		}

	// Inventory slots selection
	case ebiten.IsKeyPressed(ebiten.KeyDigit1):
		game.inventory.SelectSlot(0)
	case ebiten.IsKeyPressed(ebiten.KeyDigit2):
		game.inventory.SelectSlot(1)
	case ebiten.IsKeyPressed(ebiten.KeyDigit3):
		game.inventory.SelectSlot(2)
	case ebiten.IsKeyPressed(ebiten.KeyDigit4):
		game.inventory.SelectSlot(3)
	case ebiten.IsKeyPressed(ebiten.KeyDigit5):
		game.inventory.SelectSlot(4)
	}

	// Debug-mode only keybinds
	if config.DebugMode {
		switch {
		// Re-generate all loaded chunks on F4+G
		case ebiten.IsKeyPressed(ebiten.KeyF4) && inpututil.IsKeyJustPressed(ebiten.KeyG):
			for _, chunk := range game.world.AllLoadedChunks() {
				game.world.Generator().GenerateImmediately(chunk)
			}
		// Make the player go faaaaaaaast
		case ebiten.IsKeyPressed(ebiten.KeyF4) && inpututil.IsKeyJustPressed(ebiten.KeyS):
			game.superSpeed = !game.superSpeed

		case inpututil.IsKeyJustPressed(ebiten.KeyF6):
			game.world.SetBlock(lookingAt.X, lookingAt.Y, types.NewFurnaceBlock())
		}
	}

	_, yoff := ebiten.Wheel()
	if yoff < 0 {
		game.inventory.SelectSlot(game.inventory.SelectedSlotIndex() + 1)
	} else if yoff > 0 {
		game.inventory.SelectSlot(game.inventory.SelectedSlotIndex() - 1)
	}
}

func (game *Game) updateLogic() {
	if game.paused {
		return
	}

	game.world.Update()
	game.player.Update(game.superSpeed)
	game.compass.Update()

	// perform autosave each N ticks
	if scene_manager.Ticks()%config.WorldAutosaveDelay == 0 {
		game.Save()
	}
}

func (game *Game) handleEvents() {
	for _, ev := range event.GetEvents() {
		switch ev.Type() {
		case event.CaveEnter:
			// Move player a bit from the cave entrance, so when the world is loaded back,
			// the player won't be immediately teleported to cave
			vel := game.player.Velocity()
			game.player.Move(types.Vec2f{
				X: -vel.X * 10,
				Y: -vel.Y * 10,
			})

			// save the previous world before switching to a new one
			game.Save()

			caveID := ev.Args().(event.CaveEnteredArgs).ID

			metadata := types.Save{
				Name:      game.world.Metadata().Name,
				BaseUUID:  game.world.Metadata().BaseUUID,
				UUID:      caveID,
				Seed:      int64(caveID.ID()),
				WorldType: world_type.Cave,
				Size:      world.SizeForWorldType(world_type.Cave),
			}

			var newWorld *world.World
			// Check if cave already exists on disk
			if world.ExistsOnDisk(metadata) {
				newWorld = world.Load(metadata.BaseUUID, metadata.UUID)
			} else {
				newWorld = world.NewWorld(metadata)
			}

			game.playerStack.Push(player.NewPlayer(newWorld))
			game.player = game.playerStack.Top()

			// don't place cave exit if that chunk already exists on disk, so we don't overwrite it
			if !world.ChunkExistsOnDisk(newWorld.Metadata(), uint64(game.player.X+1)/16, uint64(game.player.Y)/16) {
				// place a cave exit next to the player
				newWorld.SetBlock(uint64(game.player.X)+1, uint64(game.player.Y), types.NewCaveExitBlock())
			}

			game.world = newWorld
			game.Save()
		case event.CaveExit:
			game.Save()
			game.playerStack.Pop()
			game.player = game.playerStack.Top()
			// reload the world
			game.world = world.Load(game.player.SelectedWorld.BaseUUID, game.player.SelectedWorld.UUID)
			game.Save()
		}
	}
}

func (game *Game) Update() {
	types.SetCurrentPlayer(game.player)
	types.SetCurrentWorld(game.world)
	game.processInput()
	game.updateLogic()
	game.handleEvents()
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.world.Render(screen, game.player.X, game.player.Y, config.UIScaling)

	if !game.inventory.SelectedSlot().Empty {
		screenPos := world.BlockToScreen(screen, types.Vec2f{X: game.player.X, Y: game.player.Y}, game.player.LookingAt(), config.UIScaling)
		tex := assets.Texture("outline1").Texture()
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(screenPos.X, screenPos.Y)
		opts.GeoM.Scale(config.UIScaling, config.UIScaling)
		screen.DrawImage(tex, opts)
	}
	game.player.Render(screen, config.UIScaling, game.paused)

	game.inventory.Render(screen)

	ui.ImmediateDraw(screen,
		ui.PositionSelf(ui.PositionTopRight, ui.Padding(0.5,
			game.compass,
		)))

	if config.DebugMode {
		font.RenderFont(screen,
			fmt.Sprintf(
				heredoc.Doc(`
					player pos:		%.2f, %.2f
					world seed:		%v
					UI scaling:		%v

					FPS:			%.0f
					TPS:			%.0f
				`),
				game.player.X, game.player.Y, game.world.Seed(), config.UIScaling, ebiten.ActualFPS(), ebiten.ActualTPS(),
			),
			0, 0, colors.C("black"),
		)

	}
}

func (game *Game) Destroy() {
	game.Save()
	log.Println("GameScene.Destroy() called")
}
