// Provides useful abstractions over base SDL2

package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Engine struct {
	Win     *sdl.Window
	Ren     *sdl.Renderer
	Surface *sdl.Surface // Main surface of the window
}

var (
	GlobalEngine *Engine = nil
)

type WindowParams struct {
	Title         string
	Width, Height int32
	Flags         uint32

	// position on the screen
	X, Y int32
}

func init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Sprintf("Failed to init sdl - %v", err))
	}

	if err := img.Init(img.INIT_PNG); err != nil {
		panic(fmt.Sprintf("Failed to init sdl_image - %v", err))
	}

	if err := ttf.Init(); err != nil {
		panic(fmt.Sprintf("Failed to init sdl_ttf - %v", err))
	}
}

func Create(wp WindowParams) (*Engine, error) {
	window, err := sdl.CreateWindow(wp.Title, wp.X, wp.Y, wp.Width, wp.Height, wp.Flags)
	if err != nil {
		return nil, fmt.Errorf("failed to create window - %v", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer - %v", err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		return nil, fmt.Errorf("familed to get window surface - %v", err)
	}

	return &Engine{
		Win:     window,
		Ren:     renderer,
		Surface: surface,
	}, nil
}

// Quits the engine completely, destroying and closing everything
func Quit(engine *Engine) {
	engine.Ren.Destroy()
	engine.Win.Destroy()
	sdl.Quit()
	img.Quit()
	ttf.Quit()
}
