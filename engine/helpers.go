/*
	Helper functions to simplify usual operations
	( font rendering, converting surfaces to textures, etc. )
*/

package engine

import (
	"strings"

	"github.com/3elDU/bamboo/engine/texture"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Clears the screen with specified color
func (e *Engine) ClearColor(r, g, b, a uint8) error {
	// saving previous color
	pr, pg, pb, pa, _ := e.Ren.GetDrawColor()
	defer e.Ren.SetDrawColor(pr, pg, pb, pa)

	e.Ren.SetDrawColor(r, g, b, a)
	e.Ren.Clear()

	return nil
}

// Clears the texture with given color, preserving previous color and render target
// Updates texture after rendering
func (e *Engine) ClearTexture(texture *sdl.Texture, r, g, b, a uint8) error {
	// saving previous color
	pr, pg, pb, pa, _ := e.Ren.GetDrawColor()
	defer e.Ren.SetDrawColor(pr, pg, pb, pa)

	prevRenderingTarget := e.Ren.GetRenderTarget()
	defer e.Ren.SetRenderTarget(prevRenderingTarget)

	e.Ren.SetRenderTarget(texture)
	e.Ren.SetDrawColor(r, g, b, a)
	e.Ren.Clear()
	e.Ren.Present()

	return nil
}

func (e *Engine) FillRectF(x, y, w, h float32, clr sdl.Color) {
	pr, pg, pb, pa, _ := e.Ren.GetDrawColor()
	defer e.Ren.SetDrawColor(pr, pg, pb, pa)

	e.Ren.SetDrawColor(clr.R, clr.G, clr.B, clr.A)
	e.Ren.FillRectF(&sdl.FRect{
		X: x, Y: y, W: w, H: h,
	})
}

// Renders texture preserving it's original width and height
// Without need to worry about creating Rect's on your own
func (e *Engine) RenderTexture(tex *sdl.Texture, x, y int32) error {
	width, height := texture.Dimensions(tex)
	err := e.Ren.Copy(tex, nil, &sdl.Rect{
		X: x, Y: y,
		W: width, H: height,
	})
	return err
}

// The same as RenderTexture(), but automatically converts surface to texture
func (e *Engine) RenderSurface(surf *sdl.Surface, x, y int32) error {
	tex, err := e.Ren.CreateTextureFromSurface(surf)
	if err != nil {
		return err
	}

	e.RenderTexture(tex, x, y)
	return nil
}

// Helper function to easily draw some text on the screen
func (e *Engine) RenderFont(font *ttf.Font, x, y int32, text string, color sdl.Color) error {
	// properly handle multi-line text
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		surf, err := font.RenderUTF8Solid(line, color)
		if err != nil {
			return err
		}

		err = e.RenderSurface(surf, x, y+int32(i*font.Height()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) CreateTexture(width, height int32) (*sdl.Texture, error) {
	format, err := e.Win.GetPixelFormat()
	if err != nil {
		return nil, err
	}

	tex, err := e.Ren.CreateTexture(format, sdl.TEXTUREACCESS_STATIC, width, height)
	if err != nil {
		return nil, err
	}

	return tex, nil
}
