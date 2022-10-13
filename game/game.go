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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	widgets      *widget.WidgetContainer
	debugWidgets *widget.WidgetContainer

	world  *world.World
	player *Player

	debugInfoVisible bool
}

func Create() *Game {
	game := &Game{
		widgets:      widget.NewWidgetContainer(),
		debugWidgets: widget.NewWidgetContainer(),

		world:  world.NewWorld(rand.Int63()),
		player: &Player{0, 0, 0, 0},

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

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.world.Render(screen, game.player.X, game.player.Y, game.debugInfoVisible)

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
			fmt.Sprint(game.world.Seed()),
			0, 24, colors.Black,
		)

		engine.RenderFont(screen, asset_loader.DefaultFont(),
			fmt.Sprint(game.debugInfoVisible),
			0, 48, colors.Black)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
