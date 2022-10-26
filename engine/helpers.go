/*
	Helper functions to simplify usual operations
*/

package engine

import (
	"image/color"
	"strings"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// RenderFontWithOptions
// arguments x and y are top-left coordinates, because RenderFontWithOptions calculates text position
// correctly
func RenderFontWithOptions(dest *ebiten.Image, face font.Face, s string, x, y float64, clr color.Color, scaling float64) {
	// TODO: Splitting by bare '\n' may not work on Windows platforms
	lines := strings.Split(s, "\n")

	// calculate a darkened tone of the main color, used later to draw the shadow
	shadowClr := colors.Complementary(clr)

	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.ColorM.ScaleWithColor(shadowClr)
	bgOpts.GeoM.Scale(scaling, scaling)
	bgOpts.GeoM.Translate(x+2*scaling, y+2*scaling)

	fgOpts := &ebiten.DrawImageOptions{}
	fgOpts.ColorM.ScaleWithColor(clr)
	fgOpts.GeoM.Scale(scaling, scaling)
	fgOpts.GeoM.Translate(x, y)

	for _, line := range lines {
		bounds := text.BoundString(face, line)
		vy := float64(bounds.Dy()) * scaling

		bgOpts.GeoM.Translate(0, float64(vy))
		fgOpts.GeoM.Translate(0, float64(vy))

		// shadow
		text.DrawWithOptions(dest, line, face, bgOpts)
		// main text
		text.DrawWithOptions(dest, line, face, fgOpts)

	}
}

// The same as RenderFontWithOptions, but with fewer options
func RenderFont(dest *ebiten.Image, s string, x, y float64, clr color.Color) {
	RenderFontWithOptions(
		dest,
		asset_loader.DefaultFont(),
		s,
		x, y,
		clr,
		1.0,
	)
}
