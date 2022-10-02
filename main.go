package main

import (
	"math/rand"
	"time"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/game"

	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	eng, err := engine.Create(engine.WindowParams{
		Title: "Hello, SDL2!",
		Width: 640, Height: 480,
		Flags: sdl.WINDOW_RESIZABLE,
	})
	eng.Ren.RenderSetVSync(true)

	if err != nil {
		panic(err)
	}

	defer engine.Quit(eng)

	game := game.Create(eng, "./assets/")
	game.Run()
}
