/*
	Helper functions for font rendering
*/

package font

import (
	"image"
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/event"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

// Do not change this unless you really know what you are doing
const (
	CharWidth  = 5
	CharHeight = 10
)

var charMap = map[rune]types.Vec2u{
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

	'N': {X: 0, Y: 10},
	'O': {X: 5, Y: 10},
	'P': {X: 10, Y: 10},
	'Q': {X: 15, Y: 10},
	'R': {X: 20, Y: 10},
	'S': {X: 25, Y: 10},
	'T': {X: 30, Y: 10},
	'U': {X: 35, Y: 10},
	'V': {X: 40, Y: 10},
	'W': {X: 45, Y: 10},
	'X': {X: 50, Y: 10},
	'Y': {X: 55, Y: 10},
	'Z': {X: 60, Y: 10},

	'a': {X: 0, Y: 20},
	'b': {X: 5, Y: 20},
	'c': {X: 10, Y: 20},
	'd': {X: 15, Y: 20},
	'e': {X: 20, Y: 20},
	'f': {X: 25, Y: 20},
	'g': {X: 30, Y: 20},
	'h': {X: 35, Y: 20},
	'i': {X: 40, Y: 20},
	'j': {X: 45, Y: 20},
	'k': {X: 50, Y: 20},
	'l': {X: 55, Y: 20},
	'm': {X: 60, Y: 20},

	'n': {X: 0, Y: 30},
	'o': {X: 5, Y: 30},
	'p': {X: 10, Y: 30},
	'q': {X: 15, Y: 30},
	'r': {X: 20, Y: 30},
	's': {X: 25, Y: 30},
	't': {X: 30, Y: 30},
	'u': {X: 35, Y: 30},
	'v': {X: 40, Y: 30},
	'w': {X: 45, Y: 30},
	'x': {X: 50, Y: 30},
	'y': {X: 55, Y: 30},
	'z': {X: 60, Y: 30},

	'.': {X: 0, Y: 40},
	',': {X: 5, Y: 40},
	';': {X: 10, Y: 40},
	'/': {X: 15, Y: 40},
	'(': {X: 20, Y: 40},
	')': {X: 25, Y: 40},
	'{': {X: 30, Y: 40},
	'}': {X: 35, Y: 40},
	'?': {X: 40, Y: 40},
	'-': {X: 45, Y: 40},
	'+': {X: 50, Y: 40},
	'%': {X: 55, Y: 40},
	'$': {X: 60, Y: 40},

	'1':  {X: 0, Y: 50},
	'2':  {X: 5, Y: 50},
	'3':  {X: 10, Y: 50},
	'4':  {X: 15, Y: 50},
	'5':  {X: 20, Y: 50},
	'6':  {X: 25, Y: 50},
	'7':  {X: 30, Y: 50},
	'8':  {X: 35, Y: 50},
	'9':  {X: 40, Y: 50},
	'0':  {X: 45, Y: 50},
	'"':  {X: 50, Y: 50},
	'\'': {X: 55, Y: 50},
	':':  {X: 60, Y: 50},

	'>': {X: 0, Y: 60},
	'<': {X: 5, Y: 60},
}

var cacheMap map[rune]*ebiten.Image

// drops the cache
func reload(_ interface{}) {
	cacheMap = make(map[rune]*ebiten.Image)
}

func init() {
	cacheMap = make(map[rune]*ebiten.Image)
	event.Subscribe(event.Reload, reload)
}

func RenderFontWithOptions(dest *ebiten.Image, s string, x, y float64, clr color.Color, scaling float64, shadow bool) {
	scaling *= config.UIScaling
	complementaryClr := colors.Complementary(clr)

	lines := strings.Split(s, "\n")

	opts := &ebiten.DrawImageOptions{}
	opts2 := &ebiten.DrawImageOptions{}

	for j, line := range lines {
		for i, char := range line {

			img, cached := cacheMap[char]
			if !cached {
				coords, exists := charMap[char]
				if !exists {
					continue
				}

				img = ebiten.NewImageFromImage(asset_loader.DefaultFont().SubImage(
					image.Rect(int(coords.X), int(coords.Y), int(coords.X+CharWidth), int(coords.Y+CharHeight))))
				cacheMap[char] = img
			}

			opts.GeoM.Reset()
			opts.GeoM.Scale(scaling, scaling)
			opts.GeoM.Translate(
				x+float64(i)*(float64(CharWidth+1)*scaling),
				y+float64(j)*(float64(CharHeight+1)*scaling),
			)
			opts.ColorScale.Reset()
			opts.ColorScale.ScaleWithColor(clr)

			if shadow {
				opts2.GeoM.Reset()
				opts2.GeoM.Scale(scaling, scaling)
				opts2.GeoM.Translate(
					x+float64(i)*(float64(CharWidth+1)*scaling)+scaling,
					y+float64(j)*(float64(CharHeight+1)*scaling)+scaling,
				)
				opts2.ColorScale.Reset()
				opts2.ColorScale.ScaleWithColor(complementaryClr)
				dest.DrawImage(img, opts2)
			}

			dest.DrawImage(img, opts)
		}
	}
}

// RenderFont is the same as RenderFontWithOptions, but with fewer options
func RenderFont(dest *ebiten.Image, s string, x, y float64, clr color.Color) {
	RenderFontWithOptions(
		dest,
		s,
		x, y,
		clr,
		1,
		true,
	)
}

// GetStringWidth Returns width of the given string in pixels
// Also handles multi-line strings properly
func GetStringWidth(s string, scaling float64) float64 {
	lines := strings.Split(s, "\n")
	if len(lines) == 1 {
		return float64(utf8.RuneCountInString(s)*(CharWidth+1)) * config.UIScaling * scaling
	} else if len(lines) > 1 {
		max := 0
		for _, line := range lines {
			runeCount := utf8.RuneCountInString(line)
			if runeCount > max {
				max = runeCount
			}
		}
		return float64(max*(CharWidth+1)) * config.UIScaling
	}

	return 0
}

// GetStringHeight returns height of the given string in pixels
// Also handles multi-line stirngs properly
func GetStringHeight(s string, scaling float64) float64 {
	nLines := len(strings.Split(s, "\n"))
	return float64(nLines*(CharHeight+1)) * config.UIScaling * scaling
}

// GetStringSize returns width and height of given string in pixels
// Also handles multi-line strings properly
func GetStringSize(s string, scaling float64) (width float64, height float64) {
	width = GetStringWidth(s, scaling)
	height = GetStringHeight(s, scaling)
	return
}
