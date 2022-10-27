package game

import (
	"fmt"
	"log"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/scene"
	"github.com/3elDU/bamboo/engine/widget"
	"github.com/3elDU/bamboo/engine/world"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameScene struct {
	widgets      *widget.WidgetContainer
	debugWidgets *widget.WidgetContainer

	paused    bool
	pauseMenu *pauseMenu

	world                  *world.World
	player                 *Player
	playerRenderingOptions *ebiten.DrawImageOptions

	scaling         float64
	scalingVelocity float64 // for smooth scaling animation

	debugInfoVisible bool
}

func New(seed int64) *GameScene {
	game := &GameScene{
		widgets:      widget.NewWidgetContainer(),
		debugWidgets: widget.NewWidgetContainer(),

		pauseMenu: newPauseMenu(),

		world:                  world.NewWorld(seed),
		player:                 &Player{0, 0, 0, 0},
		playerRenderingOptions: &ebiten.DrawImageOptions{},

		scaling: 1.0,

		debugInfoVisible: true,
	}

	/*
		game.Widgets = append(game.Widgets,
			&widgets.TextureWidget{image: asset_loader.Texture("test")},
		)
	*/

	game.debugWidgets.AddTextWidget(
		"debug",
		&widgets.PerfWidget{Color: colors.Black},
	)

	return game
}

func (game *GameScene) Update(manager *scene.SceneManager) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.paused = !game.paused
		log.Printf("Escape pressed. Toggled pause menu. (%v)", game.paused)
	}

	if !game.paused {
		// Check for key presses
		switch {
		// F3 toggles visibility of debug widgets
		case inpututil.IsKeyJustPressed(ebiten.KeyF3):
			game.debugInfoVisible = !game.debugInfoVisible
			log.Printf("Toggled visibility of debug info. (%v)", game.debugInfoVisible)
		}

		// scale the map, using scroll wheel
		_, yvel := ebiten.Wheel()
		game.scalingVelocity += yvel * 0.001

		game.player.Update(MovementVector{
			Left:  ebiten.IsKeyPressed(ebiten.KeyA),
			Right: ebiten.IsKeyPressed(ebiten.KeyD),
			Up:    ebiten.IsKeyPressed(ebiten.KeyW),
			Down:  ebiten.IsKeyPressed(ebiten.KeyS),
		})

		game.world.Update(game.player.X, game.player.Y)

		game.widgets.Update()

		if game.debugInfoVisible {
			game.debugWidgets.Update()
		}

		game.scaling += game.scalingVelocity
		game.scaling = util.Clamp(game.scaling, 1.00, 4.00)
		game.scalingVelocity *= 0.95
	} else {
		switch game.pauseMenu.ButtonPressed() {
		case continueButtonPressed:
			game.paused = false
		case exitButtonPressed:
			manager.End()
		}
	}

	return nil
}

func (game *GameScene) Draw(screen *ebiten.Image) {
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
	screen.DrawImage(asset_loader.Texture("person"), game.playerRenderingOptions)

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

		engine.RenderFont(screen,
			fmt.Sprintf("player pos %.2f, %.2f", game.player.X, game.player.Y),
			0, 0, colors.Black)

		engine.RenderFont(screen,
			fmt.Sprintf("world seed %v", game.world.Seed()),
			0, 24, colors.Black,
		)

		engine.RenderFont(screen,
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

func (g *GameScene) Destroy() {
	log.Println("GameScene.Destroy() called")
}
