package main

import (
	"math/rand"
	"time"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/game"

	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	eng, err := engine.Create(engine.WindowParams{
		Title: "bamboo devtest",
		Width: 960, Height: 640,
		Flags: sdl.WINDOW_RESIZABLE,
	})
	if err != nil {
		panic(err)
	}
	engine.GlobalEngine = eng

	eng.Ren.RenderSetVSync(true)
	defer engine.Quit(eng)

	game := game.Create(eng, "./assets/")
	asset_loader.Assets = game.Assets
	game.Run()
}
