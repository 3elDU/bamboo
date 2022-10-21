/*
	Helper functions to simplify usual operations
*/

package engine

import (
	"image/color"
	"strings"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// RenderFont
// arguments x and y are top-left coordinates, because RenderFont calculates text position
// correctly
func RenderFont(dest *ebiten.Image, face font.Face, s string, x, y int, clr color.Color) {
	// TODO: Splitting by bare '\n' may not work on Windows platforms
	lines := strings.Split(s, "\n")

	// calculate a darkened tone of the main color, used later to draw the shadow
	shadowClr := colors.Complementary(clr)

	y -= int(config.FontSize / 16)
	for _, line := range lines {
		bounds := text.BoundString(face, line)
		y += bounds.Dy() + int(config.FontSize/16)

		// draw a "shadow" under the main text
		text.Draw(dest, line, face, x+2, y+2, shadowClr)

		text.Draw(dest, line, face, x, y, clr)
	}
}
