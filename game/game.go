package game

import (
	"fmt"
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/3elDU/bamboo/scene_manager"
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

	world                  *world.World
	player                 *player.Player
	playerRenderingOptions *ebiten.DrawImageOptions

	scaling         float64
	scalingVelocity float64 // for smooth scaling animation

	blockInHand world.Item

	debugInfoVisible bool
}

func NewGameScene(gameWorld *world.World, player player.Player) *gameScene {
	game := &gameScene{
		widgets:      widget.NewWidgetContainer(),
		debugWidgets: widget.NewWidgetContainer(),

		pauseMenu: newPauseMenu(),

		world:                  gameWorld,
		player:                 &player,
		playerRenderingOptions: &ebiten.DrawImageOptions{},

		scaling: 1.0,

		blockInHand: world.NewCustomItem(asset_loader.ConnectedTexture("grass", false, false, false, false), blocks.Grass, 1),

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

		// Places block under the player
		case ebiten.IsKeyPressed(ebiten.KeyF):
			game.world.ChunkAtB(uint64(game.player.X), uint64(game.player.Y)).
				SetBlock(uint(game.player.X)%16, uint(game.player.Y)%16, blocks.GetBlockByID(game.blockInHand.Type()))

		// Pick up the block under the player
		case ebiten.IsKeyPressed(ebiten.KeyP):
			// FIXME: this is completely broken, until items dropped by blocks are implemented properly
			/*
				if block, ok := game.world.BlockAt(uint64(game.player.X), uint64(game.player.Y)).(types.DrawableBlock); ok {
					game.blockInHand = world.NewCustomItem(asset_loader.Texture(block.TextureName()), block.Type(), 1)
				}
			*/
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
	sw, sh := screen.Size()

	// draw the world
	game.world.Render(screen, game.player.X, game.player.Y, game.scaling)

	// Render the player
	game.playerRenderingOptions.GeoM.Reset()
	game.playerRenderingOptions.GeoM.Scale(game.scaling, game.scaling)
	game.playerRenderingOptions.GeoM.Translate(
		float64(sw)/2-8*game.scaling,
		float64(sh)/2-8*game.scaling,
	)
	screen.DrawImage(asset_loader.Texture("person").Texture(), game.playerRenderingOptions)

	// render block in hand
	if game.blockInHand != nil {
		opts := &ebiten.DrawImageOptions{}
		game.blockInHand.Texture()
		opts.GeoM.Scale(5, 5)
		opts.GeoM.Translate(0, 80)
		screen.DrawImage(game.blockInHand.Texture(), opts)
	}

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
