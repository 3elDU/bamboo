package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/widget"
	"github.com/3elDU/bamboo/engine/world"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	// SDL-related
	Engine  *engine.Engine
	Assets  *asset_loader.AssetList
	Widgets []widget.Widget

	World  *world.World
	Player *Player

	// Misc. variables
	running bool
	FPS     float64
}

func Create(engine *engine.Engine, assetsDirectory string) *Game {
	game := &Game{
		Engine:  engine,
		Assets:  asset_loader.LoadAssets(engine, assetsDirectory),
		Widgets: make([]widget.Widget, 0),

		World:  world.NewWorld(rand.Int63()),
		Player: &Player{0, 0, 0, 0},

		running: true,
	}

	game.Widgets = append(game.Widgets,
		// &widgets.TextureWidget{Texture: game.assets.Textures["test"]},
		&widgets.FPSWidget{
			Renderer: engine.Ren,
			Color:    colors.Red,
			Font:     game.Assets.DefaultFont(),
			FPS:      &game.FPS,
		},
	)

	return game
}

func (game *Game) Update() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			game.running = false
		case *sdl.KeyboardEvent:
			if t.State == sdl.RELEASED || t.Repeat > 0 {
				break
			}

			switch t.Keysym.Sym {
			/*
				case sdl.K_0:
					config.PERLIN_NOISE_SCALE_FACTOR += 5
				case sdl.K_9:
					config.PERLIN_NOISE_SCALE_FACTOR -= 5
			*/
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

	game.World.Update(game.Player.X, game.Player.Y)
}

func (game *Game) Render() {
	game.Engine.Ren.Clear()

	game.World.Render(game.Player.X, game.Player.Y)
	widget.RenderMultiple(game.Engine, game.Widgets)
	game.Engine.RenderFont(
		game.Assets.DefaultFont(), 0, 0,
		fmt.Sprintf("%v, %v", game.Player.X, game.Player.Y),
		colors.Black,
	)
	game.Engine.RenderFont(
		game.Assets.DefaultFont(), 0, 32,
		fmt.Sprint(config.PERLIN_NOISE_SCALE_FACTOR),
		colors.Black,
	)

	game.Engine.Ren.Present()
}

func (game *Game) Running() bool {
	return game.running
}

func (game *Game) Run() {
	for game.Running() {
		frameStart := time.Now()

		game.Update()
		game.Render()

		frameEnd := time.Now()

		game.FPS = 1 / frameEnd.Sub(frameStart).Seconds()
	}
}
