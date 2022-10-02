package widgets

import (
	"fmt"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/widgets"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TextureWidget struct {
	Texture *sdl.Texture
	a       widgets.Anchor
}

func (w *TextureWidget) Anchor() widgets.Anchor {
	w.a++
	if w.a > widgets.BottomRight {
		w.a = 0
	}

	return w.a
}

func (w *TextureWidget) Render() *sdl.Texture {
	return w.Texture
}

type FPSWidget struct {
	Anc    widgets.Anchor
	Color  sdl.Color
	Font   *ttf.Font
	Engine *engine.Engine
}

func (w *FPSWidget) Anchor() widgets.Anchor {
	return w.Anc
}

func (w *FPSWidget) Render() *sdl.Texture {
	surf, _ := w.Font.RenderUTF8Solid(fmt.Sprint(w.Engine.FPS()), w.Color)
	tex, _ := w.Engine.Ren.CreateTextureFromSurface(surf)
	return tex
}
