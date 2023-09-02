package colors

import (
	"image/color"
	"strings"

	"github.com/teacat/noire"
)

var Colors map[string]color.Color = map[string]color.Color{
	"black": color.Black,
	"white": color.White,

	"lightgray": color.RGBA{R: 0xC0, G: 0xCB, B: 0xDC, A: 0xFF},
	"darkgray":  color.RGBA{R: 0x3F, G: 0x3F, B: 0x3F, A: 0xFF},

	"red":        color.RGBA{R: 0xE4, G: 0x3B, B: 0x44, A: 0xFF},
	"green":      color.RGBA{R: 0x63, G: 0xC7, B: 0x4D, A: 0xFF},
	"blue":       color.RGBA{R: 0x00, G: 0x99, B: 0xDB, A: 0xFF},
	"darkblue":   color.RGBA{R: 0x12, G: 0x4E, B: 0x89, A: 0xFF},
	"yellow":     color.RGBA{R: 0xFE, G: 0xE7, B: 0x61, A: 0xFF},
	"cyan":       color.RGBA{R: 0x2C, G: 0xE8, B: 0xF5, A: 0xFF},
	"orange":     color.RGBA{R: 0xFE, G: 0xAE, B: 0x34, A: 0xFF},
	"darkorange": color.RGBA{R: 0xF7, G: 0x76, B: 0x22, A: 0xFF},
	"darkviolet": color.RGBA{R: 0x68, G: 0x38, B: 0x6c, A: 0xFF},
	"violet":     color.RGBA{R: 0xB5, G: 0x50, B: 0x88, A: 0xFF},

	"darkgreen1": color.RGBA{R: 0x3E, G: 0x89, B: 0x48, A: 0xFF},
	"darkgreen2": color.RGBA{R: 0x26, G: 0x5C, B: 0x42, A: 0xFF},
	"darkgreen3": color.RGBA{R: 0x19, G: 0x3C, B: 0x3E, A: 0xFF},
}

// Get a color by the name. Casing does not matter
func C(name string) color.Color {
	if clr, exists := Colors[strings.ToLower(name)]; exists {
		return clr
	} else {
		return color.RGBA{R: 255, G: 0, B: 255, A: 255}
	}
}

func Complementary(clr color.Color) color.Color {
	switch clr {
	case C("black"):
		return C("lightgray")
	case C("white"):
		return C("darkgray")

	case C("darkgreen1"):
		return C("darkgreen3")

	default:
		r, g, b, _ := clr.RGBA()
		origColor := noire.NewRGB(float64(r)/255, float64(g)/255, float64(b)/255)

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
