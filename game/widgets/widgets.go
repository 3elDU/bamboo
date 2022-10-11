package widgets

import (
	"fmt"

	"github.com/3elDU/bamboo/engine/widget"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TextureWidget struct {
	Texture *sdl.Texture
	anchor  widget.Anchor
}

func (w *TextureWidget) Anchor() widget.Anchor {
	w.anchor++
	if w.anchor > widget.BottomRight {
		w.anchor = 0
	}

	return w.anchor
}

func (w *TextureWidget) Render() *sdl.Texture {
	return w.Texture
}

type FPSWidget struct {
	Renderer *sdl.Renderer
	Color    sdl.Color
	Font     *ttf.Font

	// FIXME: Won't work in multi-threaded envirionment
	FPS *float64
}

func (w *FPSWidget) Anchor() widget.Anchor {
	return widget.TopRight
}

func (w *FPSWidget) Render() *sdl.Texture {
	if w.FPS == nil {
		panic("Invalid pointer to FPS!")
	}

	surf, _ := w.Font.RenderUTF8Solid(fmt.Sprint(int(*w.FPS)), w.Color)
	defer surf.Free()
	tex, _ := w.Renderer.CreateTextureFromSurface(surf)
	return tex
}
