package colors

import (
	"image/color"

	"github.com/teacat/noire"
)

var (
	Black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	White = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	Gray  = color.RGBA{R: 0xC0, G: 0xCB, B: 0xDC, A: 0xFF}

	Red        = color.RGBA{R: 0xE4, G: 0x3B, B: 0x44, A: 0xFF}
	Green      = color.RGBA{R: 0x63, G: 0xC7, B: 0x4D, A: 0xFF}
	Blue       = color.RGBA{R: 0x00, G: 0x99, B: 0xDB, A: 0xFF}
	DarkBlue   = color.RGBA{R: 0x12, G: 0x4E, B: 0x89, A: 0xFF}
	Yellow     = color.RGBA{R: 0xFE, G: 0xE7, B: 0x61, A: 0xFF}
	Cyan       = color.RGBA{R: 0x2C, G: 0xE8, B: 0xF5, A: 0xFF}
	Orange     = color.RGBA{R: 0xFE, G: 0xAE, B: 0x34, A: 0xFF}
	DarkOrange = color.RGBA{R: 0xF7, G: 0x76, B: 0x22, A: 0xFF}

	DarkGreen1 = color.RGBA{R: 0x3E, G: 0x89, B: 0x48, A: 0xFF}
	DarkGreen2 = color.RGBA{R: 0x26, G: 0x5C, B: 0x42, A: 0xFF}
	DarkGreen3 = color.RGBA{R: 0x19, G: 0x3C, B: 0x3E, A: 0xFF}
)

func Complementary(clr color.Color) color.Color {
	switch clr {
	case Black:
		return Gray
	case White:
		return Black

	case Red:
		return Black
	case Yellow:
		return DarkOrange
	case DarkOrange:
		return Yellow

	case Green, DarkGreen1:
		return DarkGreen3
	case DarkGreen2, DarkGreen3:
		return Green

	case Cyan, Blue:
		return DarkBlue
	case DarkBlue:
		return Cyan

	default:
		r, g, b, _ := clr.RGBA()
		origColor := noire.NewRGB(float64(r)/256, float64(g)/256, float64(b)/256)

		var newColor noire.Color
		if origColor.IsDark() {
			newColor = origColor.Lighten(0.4)
		} else {
			newColor = origColor.Darken(0.4)
		}

		nr, ng, nb := newColor.RGB()

		return color.RGBA{R: uint8(nr * 255), G: uint8(ng * 255), B: uint8(nb * 255), A: 255}
	}
}
