package game

import (
	"math/rand"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	engine_widgets "github.com/3elDU/bamboo/engine/widgets"
	"github.com/3elDU/bamboo/engine/worldgen"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/veandco/go-sdl2/sdl"
)

type gameState struct {
	engine  *engine.Engine
	world   *worldgen.World
	widgets []engine_widgets.Widget
	running bool
	assets  *asset_loader.AssetList
}

func Create(engine *engine.Engine, assetsDirectory string) *gameState {
	game := &gameState{
		engine:  engine,
		widgets: make([]engine_widgets.Widget, 0),
		running: true,
		assets:  asset_loader.LoadAssets(engine, assetsDirectory),
		world:   worldgen.New(rand.Int63(), config.PERLIN_NOISE_SCALE_FACTOR),
	}

	game.widgets = append(game.widgets,
		// &widgets.TextureWidget{Texture: game.assets.Textures["test"]},
		&widgets.FPSWidget{
			Anc:    engine_widgets.TopRight,
			Color:  sdl.Color{R: 255, G: 0, B: 0, A: 255},
			Font:   game.assets.Fonts["font"],
			Engine: engine,
		},
	)

	return game
}

func (game *gameState) Update() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.GetType() {
		case sdl.QUIT:
			game.running = false
		}
	}
}

func (game *gameState) Render() {
	game.engine.Ren.Clear()

	w, h := game.engine.Win.GetSize()
	const pixelSize = 32
	for x := int32(0); x < w/pixelSize; x++ {
		for y := int32(0); y < h/pixelSize; y++ {
			screenX, screenY := x*pixelSize, y*pixelSize

			clr := game.world.Block(float64(x), float64(y))

			game.engine.Ren.SetDrawColor(clr.R, clr.G, clr.B, clr.A)
			game.engine.Ren.FillRect(&sdl.Rect{
				X: screenX, Y: screenY,
				W: pixelSize, H: pixelSize,
			})
		}
	}

	engine_widgets.RenderMultiple(game.engine, game.widgets)

	game.engine.Ren.Present()
}

func (game *gameState) Running() bool {
	return game.running
}

func (game *gameState) Run() {
	game.engine.Run(game)
}
