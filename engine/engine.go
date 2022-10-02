package engine

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Engine struct {
	Win     *sdl.Window
	Ren     *sdl.Renderer
	Surface *sdl.Surface // Surface of the screen

	fps int
}

// Render() and Update() functions will be called each frame
// It is up to user, what to include in these functions
// Note that Update is called before Render
type Game interface {
	Render()
	Update()

	// Engine calls this function each frame, to know if it should execute further
	Running() bool
}

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
		window, renderer, surface, 0,
	}, nil
}

func (e *Engine) FPS() int {
	return e.fps
}

// The main loop! Runs the game until something calls engine.Quit()
func (e *Engine) Run(g Game) {
	for g.Running() {
		frameStart := time.Now()

		g.Update()
		g.Render()

		frameEnd := time.Now()

		e.fps = int(1 / frameEnd.Sub(frameStart).Seconds())
	}
}

// Quits the engine completely, destroying and closing everything
func Quit(engine *Engine) {
	engine.Ren.Destroy()
	engine.Win.Destroy()
	sdl.Quit()
	img.Quit()
}
