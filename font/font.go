/*
	Helper functions for font rendering
*/

package font

import (
	"image"
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

// Do not change this unless you really know what you are doing
const (
	FontWidth  = 5
	FontHeight = 7
)

var CharMap map[rune]types.Coords2u = map[rune]types.Coords2u{
	'A': {X: 0, Y: 0},
	'B': {X: 5, Y: 0},
	'C': {X: 10, Y: 0},
	'D': {X: 15, Y: 0},
	'E': {X: 20, Y: 0},
	'F': {X: 25, Y: 0},
	'G': {X: 30, Y: 0},
	'H': {X: 35, Y: 0},
	'I': {X: 40, Y: 0},
	'J': {X: 45, Y: 0},
	'K': {X: 50, Y: 0},
	'L': {X: 55, Y: 0},
	'M': {X: 60, Y: 0},

	'N': {X: 0, Y: 7},
	'O': {X: 5, Y: 7},
	'P': {X: 10, Y: 7},
	'Q': {X: 15, Y: 7},
	'R': {X: 20, Y: 7},
	'S': {X: 25, Y: 7},
	'T': {X: 30, Y: 7},
	'U': {X: 35, Y: 7},
	'V': {X: 40, Y: 7},
	'W': {X: 45, Y: 7},
	'X': {X: 50, Y: 7},
	'Y': {X: 55, Y: 7},
	'Z': {X: 60, Y: 7},

	'a': {X: 0, Y: 14},
	'b': {X: 5, Y: 14},
	'c': {X: 10, Y: 14},
	'd': {X: 15, Y: 14},
	'e': {X: 20, Y: 14},
	'f': {X: 25, Y: 14},
	'g': {X: 30, Y: 14},
	'h': {X: 35, Y: 14},
	'i': {X: 40, Y: 14},
	'j': {X: 45, Y: 14},
	'k': {X: 50, Y: 14},
	'l': {X: 55, Y: 14},
	'm': {X: 60, Y: 14},

	'n': {X: 0, Y: 21},
	'o': {X: 5, Y: 21},
	'p': {X: 10, Y: 21},
	'q': {X: 15, Y: 21},
	'r': {X: 20, Y: 21},
	's': {X: 25, Y: 21},
	't': {X: 30, Y: 21},
	'u': {X: 35, Y: 21},
	'v': {X: 40, Y: 21},
	'w': {X: 45, Y: 21},
	'x': {X: 50, Y: 21},
	'y': {X: 55, Y: 21},
	'z': {X: 60, Y: 21},

	'.': {X: 0, Y: 28},
	',': {X: 5, Y: 28},
	';': {X: 10, Y: 28},
	'/': {X: 15, Y: 28},
	'(': {X: 20, Y: 28},
	')': {X: 25, Y: 28},
	'{': {X: 30, Y: 28},
	'}': {X: 35, Y: 28},
	'?': {X: 40, Y: 28},
	'-': {X: 45, Y: 28},
	'+': {X: 50, Y: 28},
	'%': {X: 55, Y: 28},
	'$': {X: 60, Y: 28},

	'1':  {X: 0, Y: 35},
	'2':  {X: 5, Y: 35},
	'3':  {X: 10, Y: 35},
	'4':  {X: 15, Y: 35},
	'5':  {X: 20, Y: 35},
	'6':  {X: 25, Y: 35},
	'7':  {X: 30, Y: 35},
	'8':  {X: 35, Y: 35},
	'9':  {X: 40, Y: 35},
	'0':  {X: 45, Y: 35},
	'"':  {X: 50, Y: 35},
	'\'': {X: 55, Y: 35},
	':':  {X: 60, Y: 35},
}

var cacheMap map[rune]*ebiten.Image

func init() {
	cacheMap = make(map[rune]*ebiten.Image)
}

func RenderFontWithOptions(dest *ebiten.Image, face *ebiten.Image, s string, x, y float64, clr color.Color, scaling float64) {
	scaling *= float64(config.UIScaling)

	lines := strings.Split(s, "\n")

	opts := &ebiten.DrawImageOptions{}
	opts2 := &ebiten.DrawImageOptions{}

	for j, line := range lines {
		for i, char := range line {

			img, exists := cacheMap[char]

			if !exists {
				coords, exists := CharMap[char]
				if !exists {
					// log.Panicf("char doesn't exist: %v", char)
					continue
				}

				img = ebiten.NewImageFromImage(asset_loader.DefaultFont().SubImage(
					image.Rect(int(coords.X), int(coords.Y), int(coords.X+5), int(coords.Y+7))))
				cacheMap[char] = img
			}

			opts.GeoM.Reset()
			opts.GeoM.Scale(scaling, scaling)
			opts.GeoM.Translate(
				x+float64(i)*(float64(6)*scaling),
				y+float64(j)*(float64(8)*scaling),
			)
			opts.ColorM.Reset()
			opts.ColorM.ScaleWithColor(clr)

			opts2.GeoM.Reset()
			opts2.GeoM.Scale(scaling, scaling)
			opts2.GeoM.Translate(
				x+float64(i)*(float64(6)*scaling)+scaling,
				y+float64(j)*(float64(8)*scaling)+scaling,
			)
			opts2.ColorM.Reset()
			opts2.ColorM.ScaleWithColor(colors.Complementary(clr))

			dest.DrawImage(img, opts2)
			dest.DrawImage(img, opts)
		}
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
		1,
	)
}

// Returns width of the given string in pixels
// Handles multi-line strings properly
func GetStringWidth(s string, scaling float64) int {
	lines := strings.Split(s, "\n")
	if len(lines) == 1 {
		return int(float64(utf8.RuneCountInString(s)*(FontWidth+1)*config.UIScaling) * scaling)
	} else if len(lines) > 1 {
		max := 0
		for _, line := range lines {
			runeCount := utf8.RuneCountInString(line)
			if runeCount > max {
				max = runeCount
			}
		}
		return max * (FontWidth + 1) * int(config.UIScaling)
	}

	return 0
}

// Returns height of the given string in pixels
// Handles multi-line stirngs properly
func GetStringHeight(s string, scaling float64) int {
	nLines := len(strings.Split(s, "\n"))
	return int(float64(nLines*(FontHeight+1)*int(config.UIScaling)) * scaling)
}

// Returns width and height of given string in pixels
// Handles multi-line strings properly
func GetStringSize(s string, scaling float64) (width int, height int) {
	width = GetStringWidth(s, scaling)
	height = GetStringHeight(s, scaling)
	return
}
