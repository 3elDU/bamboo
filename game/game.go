package game

import (
	"fmt"
	"math/rand"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/widget"
	"github.com/3elDU/bamboo/engine/world"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	widgets      *widget.WidgetContainer
	debugWidgets *widget.WidgetContainer

	world                  *world.World
	player                 *Player
	playerRenderingOptions *ebiten.DrawImageOptions

	scaling         float64
	scalingVelocity float64 // for smooth scaling animation

	debugInfoVisible bool
}

func Create() *Game {
	game := &Game{
		widgets:      widget.NewWidgetContainer(),
		debugWidgets: widget.NewWidgetContainer(),

		world:                  world.NewWorld(rand.Int63()),
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
		&widgets.PerfWidget{Color: colors.Black, Face: asset_loader.DefaultFont()},
	)

	return game
}

func (game *Game) Update() error {
	// Check for key presses
	switch {
	// F3 toggles visibility of debug widgets
	case inpututil.IsKeyJustPressed(ebiten.KeyF3):
		game.debugInfoVisible = !game.debugInfoVisible

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

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	sw, sh := screen.Size()

	game.world.Render(screen, game.player.X, game.player.Y, game.scaling)

	// Render the player
	game.playerRenderingOptions.GeoM.Reset()
	game.playerRenderingOptions.GeoM.Scale(game.scaling, game.scaling)
	game.playerRenderingOptions.GeoM.Translate(
		float64(sw)/2-8*game.scaling,
		float64(sh)/2-8*game.scaling,
	)
	screen.DrawImage(asset_loader.Texture("person"), game.playerRenderingOptions)

	game.widgets.Render(screen)
	if game.debugInfoVisible {
		game.debugWidgets.Render(screen)
	}

	if game.debugInfoVisible {
		// TODO: extract this to separate widgets
		// But that would require lots of architecture changed
		// Because currently, there is no way to pass custom data to a widget

		engine.RenderFont(screen, asset_loader.DefaultFont(),
			fmt.Sprintf("player pos %.2f, %.2f", game.player.X, game.player.Y),
			0, 0, colors.Black)

		engine.RenderFont(screen, asset_loader.DefaultFont(),
			fmt.Sprintf("world seed %v", game.world.Seed()),
			0, 24, colors.Black,
		)

		engine.RenderFont(screen, asset_loader.DefaultFont(),
			fmt.Sprintf("scaling %v", util.LimitFloatPrecision(game.scaling, 2)),
			0, 48, colors.Black)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
