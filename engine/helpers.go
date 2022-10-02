/*
	Helper functions to simplify usual operations
	( font rendering, converting surfaces to textures, etc. )
*/

package engine

import (
	"github.com/3elDU/bamboo/engine/texture"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

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
	surf, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return err
	}

	err = e.RenderSurface(surf, x, y)
	if err != nil {
		return err
	}

	return nil
}
