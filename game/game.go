package game

import (
	"fmt"
	"log"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/game/inventory"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/3elDU/bamboo/items"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/util"
	"github.com/3elDU/bamboo/widget"
	"github.com/3elDU/bamboo/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type gameScene struct {
	widgets      *widget.WidgetContainer
	debugWidgets *widget.WidgetContainer

	paused    bool
	pauseMenu *pauseMenu

	world     *world.World
	player    *player.Player
	inventory *inventory.Inventory

	scaling         float64
	scalingVelocity float64 // for smooth scaling animation

	debugInfoVisible bool
}

func NewGameScene(gameWorld *world.World, player player.Player) *gameScene {
	game := &gameScene{
		widgets:      widget.NewWidgetContainer(),
		debugWidgets: widget.NewWidgetContainer(),

		pauseMenu: newPauseMenu(),

		world:     gameWorld,
		player:    &player,
		inventory: inventory.NewInventory(),

		scaling: 2.0,

		debugInfoVisible: true,
	}

	game.debugWidgets.AddTextWidget(
		"debug",
		&widgets.PerfWidget{Color: colors.Black},
	)

	// perform a save immediately after the scene creation
	game.Save()

	return game
}

func (game *gameScene) Save() {
	game.world.Save()
	game.player.Save(game.world.Metadata.UUID)
}

func (game *gameScene) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.paused = !game.paused
		log.Printf("Escape pressed. Toggled pause menu. (%v)", game.paused)

		// trigger a world save, when entering pause menu
		if game.paused {
			game.Save()
		}
	}

	if !game.paused {
		// Check for key presses
		switch {
		// F3 toggles visibility of debug widgets
		case inpututil.IsKeyJustPressed(ebiten.KeyF3):
			game.debugInfoVisible = !game.debugInfoVisible
			log.Printf("Toggled visibility of debug info. (%v)", game.debugInfoVisible)

		// Interact with the nearby block
		case ebiten.IsKeyPressed(ebiten.KeyC):
			block := game.world.BlockAt(uint64(game.player.X), uint64(game.player.Y))
			drawable, ok := block.(types.DrawableBlock)
			if !ok {
				break
			}

			item := items.NewItemFromBlock(drawable)
			game.inventory.AddItem(item)

		// Use the item in hand
		case ebiten.IsKeyPressed(ebiten.KeyF):
			itemInHand := game.inventory.Slots[game.inventory.SelectedSlot].Item
			if itemInHand == nil {
				break
			}
			itemInHand.Use(game.world, types.Coords2u{
				X: uint64(game.player.X),
				Y: uint64(game.player.Y),
			})

		// Inventory slots selection
		case ebiten.IsKeyPressed(ebiten.Key1):
			game.inventory.SelectSlot(0)
		case ebiten.IsKeyPressed(ebiten.Key2):
			game.inventory.SelectSlot(1)
		case ebiten.IsKeyPressed(ebiten.Key3):
			game.inventory.SelectSlot(2)
		case ebiten.IsKeyPressed(ebiten.Key4):
			game.inventory.SelectSlot(3)
		case ebiten.IsKeyPressed(ebiten.Key5):
			game.inventory.SelectSlot(4)

		}

		// scale the map, using scroll wheel
		_, yvel := ebiten.Wheel()
		game.scalingVelocity += yvel * 0.004

		game.player.Update(player.MovementVector{
			Left:  ebiten.IsKeyPressed(ebiten.KeyA),
			Right: ebiten.IsKeyPressed(ebiten.KeyD),
			Up:    ebiten.IsKeyPressed(ebiten.KeyW),
			Down:  ebiten.IsKeyPressed(ebiten.KeyS),
		}, game.world)

		game.world.Update()

		game.widgets.Update()

		if game.debugInfoVisible {
			game.debugWidgets.Update()
		}

		game.scaling += game.scalingVelocity
		game.scaling = util.Clamp(game.scaling, 1.00, 6.00)
		game.scalingVelocity *= 0.95
	} else {
		switch game.pauseMenu.ButtonPressed() {
		case continueButtonPressed:
			game.paused = false
		case exitButtonPressed:
			game.Save()
			scene_manager.End()
		}
	}

	// perform autosave each N ticks
	if scene_manager.Ticks()%config.WorldAutosaveDelay == 0 {
		game.Save()
	}
}

func (game *gameScene) Draw(screen *ebiten.Image) {
	game.world.Render(screen, game.player.X, game.player.Y, game.scaling)
	game.player.Render(screen, game.scaling)
	game.inventory.Render(screen)

	// draw widgets
	game.widgets.Render(screen)
	if game.debugInfoVisible {
		game.debugWidgets.Render(screen)
	}

	// draw debug info
	if game.debugInfoVisible {
		// TODO: extract this to separate widgets
		// But that would require lots of architecture changed
		// Because currently, there is no way to pass custom data to a widget

		font.RenderFont(screen,
			fmt.Sprintf("player pos %.2f, %.2f", game.player.X, game.player.Y),
			0, 0, colors.Black)

		font.RenderFont(screen,
			fmt.Sprintf("world seed %v", game.world.Seed()),
			0, 24, colors.Black,
		)

		font.RenderFont(screen,
			fmt.Sprintf("scaling %v", util.LimitFloatPrecision(game.scaling, 2)),
			0, 48, colors.Black)
	}

	// draw pause menu
	if game.paused {
		err := game.pauseMenu.Draw(screen)
		if err != nil {
			log.Panicf("error while rendering pause menu - %v", err)
		}
	}
}

func (game *gameScene) Destroy() {
	game.Save()
	log.Println("GameScene.Destroy() called")
}
