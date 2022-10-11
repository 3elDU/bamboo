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
	Assets      *asset_loader.AssetList
	Widgets     []widget.Widget
	TextWidgets []widget.TextWidget

	World  *world.World
	Player *Player
}

func Create(assetsDirectory string) *Game {
	game := &Game{
		Assets:      asset_loader.LoadAssets(assetsDirectory),
		Widgets:     make([]widget.Widget, 0),
		TextWidgets: make([]widget.TextWidget, 0),

		World:  world.NewWorld(rand.Int63()),
		Player: &Player{0, 0, 0, 0},
	}

	game.TextWidgets = append(game.TextWidgets,
		// &widgets.TextureWidget{Texture: game.assets.Textures["test"]},
		&widgets.PerfWidget{
			Face:  game.Assets.DefaultFont(),
			Color: colors.Black,
		},
	)

	return game
}

func (game *Game) Update() error {
	/*
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				game.running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.RELEASED || t.Repeat > 0 {
					break
				}

				switch t.Keysym.Sym {
				case sdl.K_0:
					config.PERLIN_NOISE_SCALE_FACTOR += 5
				case sdl.K_9:
					config.PERLIN_NOISE_SCALE_FACTOR -= 5
				}
			}
		}

		keysPressed := sdl.GetKeyboardState()
		game.Player.Update(MovementVector{
			Left:  keysPressed[sdl.SCANCODE_A] == 1,
			Right: keysPressed[sdl.SCANCODE_D] == 1,
			Up:    keysPressed[sdl.SCANCODE_W] == 1,
			Down:  keysPressed[sdl.SCANCODE_S] == 1,
		})
	*/

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

	engine.RenderFont(screen, game.Assets.DefaultFont(),
		fmt.Sprintf("%v, %v", game.Player.X, game.Player.Y),
		0, 0, colors.Black)

	engine.RenderFont(screen, game.Assets.DefaultFont(),
		fmt.Sprint(config.PERLIN_NOISE_SCALE_FACTOR),
		0, 32, colors.Black,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
