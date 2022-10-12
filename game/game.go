package game

import (
	"fmt"
	"math/rand"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/widget"
	"github.com/3elDU/bamboo/engine/world"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Widgets     []widget.Widget
	TextWidgets []widget.TextWidget

	World  *world.World
	Player *Player
}

func Create() *Game {
	game := &Game{
		Widgets:     make([]widget.Widget, 0),
		TextWidgets: make([]widget.TextWidget, 0),

		World:  world.NewWorld(rand.Int63()),
		Player: &Player{0, 0, 0, 0},
	}

	/*
		game.Widgets = append(game.Widgets,
			&widgets.TextureWidget{Image: asset_loader.Texture("test")},
		)
	*/

	game.TextWidgets = append(game.TextWidgets,
		&widgets.PerfWidget{
			Face:  asset_loader.DefaultFont(),
			Color: colors.Black,
		},
	)

	return game
}

func (game *Game) Update() error {
	game.Player.Update(MovementVector{
		Left:  ebiten.IsKeyPressed(ebiten.KeyA),
		Right: ebiten.IsKeyPressed(ebiten.KeyD),
		Up:    ebiten.IsKeyPressed(ebiten.KeyW),
		Down:  ebiten.IsKeyPressed(ebiten.KeyS),
	})

	game.World.Update(game.Player.X, game.Player.Y)

	for _, w := range game.Widgets {
		w.Update()
	}

	for _, w := range game.TextWidgets {
		w.Update()
	}

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.World.Render(screen, game.Player.X, game.Player.Y)

	for _, w := range game.Widgets {
		widget.RenderWidget(screen, w)
	}

	for _, w := range game.TextWidgets {
		widget.RenderTextWidget(screen, w)
	}

	engine.RenderFont(screen, asset_loader.DefaultFont(),
		fmt.Sprintf("%v, %v", game.Player.X, game.Player.Y),
		0, 0, colors.Black)

	engine.RenderFont(screen, asset_loader.DefaultFont(),
		fmt.Sprint(config.PerlinNoiseScaleFactor),
		0, 32, colors.Black,
	)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
