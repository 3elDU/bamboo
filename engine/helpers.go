/*
	Helper functions to simplify usual operations
*/

package engine

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// RenderFont
// arguments x and y are top-left coordinates, because RenderFont calculates text position
// correctly
func RenderFont(dest *ebiten.Image, face font.Face, s string, x, y int, color color.Color) {
	// TODO: Splitting by bare '\n' may not work on Windows platforms
	lines := strings.Split(s, "\n")

	for _, line := range lines {
		bounds := text.BoundString(face, line)
		y += bounds.Dy()
		text.Draw(dest, line, face, x, y, color)
	}
}
