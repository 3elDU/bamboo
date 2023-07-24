package ui

import (
	"image"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

var tooltip TooltipTexture
var tooltip_button TooltipTexture
var tooltip_button_hover TooltipTexture
var tooltip_input TooltipTexture
var tooltip_input_focused TooltipTexture

type TooltipTexture struct {
	TopLeft  *ebiten.Image
	Top      *ebiten.Image
	TopRight *ebiten.Image

	Left   *ebiten.Image
	Center *ebiten.Image
	Right  *ebiten.Image

	BottomLeft  *ebiten.Image
	Bottom      *ebiten.Image
	BottomRight *ebiten.Image
}

// Divide the tooltip texture into sub-textures for rendering
func AssembleTooltipTexture(texture *ebiten.Image) TooltipTexture {
	topLeft := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 3, Y: 3},
	}))
	topRight := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 4, Y: 0},
		Max: image.Point{X: 7, Y: 3},
	}))
	bottomLeft := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 4},
		Max: image.Point{X: 3, Y: 7},
	}))
	bottomRight := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 4, Y: 4},
		Max: image.Point{X: 7, Y: 7},
	}))

	left := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 3},
		Max: image.Point{X: 3, Y: 4},
	}))
	right := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 4, Y: 3},
		Max: image.Point{X: 7, Y: 4},
	}))
	top := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 3, Y: 0},
		Max: image.Point{X: 4, Y: 3},
	}))
	bottom := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 3, Y: 4},
		Max: image.Point{X: 4, Y: 7},
	}))

	center := ebiten.NewImageFromImage(texture.SubImage(image.Rectangle{
		Min: image.Point{X: 3, Y: 3},
		Max: image.Point{X: 4, Y: 4},
	}))

	return TooltipTexture{
		TopLeft:     topLeft,
		Top:         top,
		TopRight:    topRight,
		Left:        left,
		Center:      center,
		Right:       right,
		BottomLeft:  bottomLeft,
		Bottom:      bottom,
		BottomRight: bottomRight,
	}
}

func init() {
	tooltip = AssembleTooltipTexture(asset_loader.Texture("tooltip").Texture())
	tooltip_button = AssembleTooltipTexture(asset_loader.Texture("tooltip_button").Texture())
	tooltip_button_hover = AssembleTooltipTexture(asset_loader.Texture("tooltip_button_hover").Texture())
	tooltip_input = AssembleTooltipTexture(asset_loader.Texture("tooltip_input").Texture())
	tooltip_input_focused = AssembleTooltipTexture(asset_loader.Texture("tooltip_input_focused").Texture())
}

// Which side of the cursor to prefer for displaying the tooltip.
// By default, the tooltip is displayed where it fits. If there are multiple sides it can fit on,
// we pick the direction with the most free pixels (pixels to the edge of the screen) in both axis.
// The default order is: BottomRight, BottomLeft, TopRight, TopLeft
type TooltipSide int

const (
	None TooltipSide = iota
	BottomRight
	BottomLeft
	TopRight
	TopLeft
)

// Checks if all tooltip corners fit inside the screen
func checkCorners(sw, sh, x, y, w, h int) bool {
	return x >= 0 && x <= sw && y >= 0 && y <= sh && // Top-left
		x+w >= 0 && x+w <= sw && // Right side
		y+h >= 0 && y+h <= sh // Bottom side
}

// Finds optimal position for tooltip on the screen.
// Returns X and Y coordinates of top-left corner.
func positionTooltip(screen *ebiten.Image, cursorX, cursorY int, width, height float64, preferredSide TooltipSide) (float64, float64) {
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()

	sides := map[TooltipSide]types.Vec2i{
		BottomRight: {X: cursorX + int(3*config.UIScaling), Y: cursorY + int(3*config.UIScaling)},
		BottomLeft:  {X: cursorX - int(width) - int(3*config.UIScaling), Y: cursorY + int(3*config.UIScaling)},
		TopRight:    {X: cursorX, Y: cursorY - int(height) - int(6*config.UIScaling)},
		TopLeft:     {X: cursorX - int(width) - int(6*config.UIScaling), Y: cursorY - int(height) - int(6*config.UIScaling)},
	}
	fitsIn := make(map[TooltipSide]bool)

	for side := range sides {
		fitsIn[side] = checkCorners(sw, sh, sides[side].X, sides[side].Y, int(width), int(height))
	}

	if preferredSide != None && fitsIn[preferredSide] {
		return float64(sides[preferredSide].X), float64(sides[preferredSide].Y)
	}

	// loop over all the sides, and return the first one that fits
	// do not use the `range` iterator, because maps in Go are unordered
	for i := BottomRight; i <= 4; i++ {
		if fitsIn[i] {
			return float64(sides[i].X), float64(sides[i].Y)
		}
	}

	// if none of the sides fit, return the preferred side.
	// or, if it isn't specified, just return the BottomRight one
	if preferredSide != None {
		return float64(sides[preferredSide].X), float64(sides[preferredSide].Y)
	} else {
		return float64(sides[BottomRight].X), float64(sides[BottomRight].Y)
	}
}

// Renders a tooltip at specified coordinates, with a specified side relative to the cursor.
// Or, if no preferredSide is specified, an optimal side will be picked.
func DrawTextTooltip(screen *ebiten.Image, cursorX, cursorY int, preferredSide TooltipSide, text string) {
	w, h := font.GetStringSize(text, 1)
	x, y := positionTooltip(screen, cursorX, cursorY, w, h, preferredSide)

	DrawTooltipBackground(screen, x, y, w, h)

	font.RenderFont(screen, text, x+3*config.UIScaling, y+3*config.UIScaling, colors.C("white"))
}

func draw_tooltip(screen *ebiten.Image, texture TooltipTexture, x, y, w, h float64) {
	opts := &ebiten.DrawImageOptions{}

	// Top left corner
	opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	opts.GeoM.Translate(x, y)
	screen.DrawImage(texture.TopLeft, opts)

	// Top right corner
	opts.GeoM.Translate(w+3*config.UIScaling, 0)
	screen.DrawImage(texture.TopRight, opts)

	// Bottom left corner
	opts.GeoM.Reset()
	opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	opts.GeoM.Translate(x, y+h+3*config.UIScaling)
	screen.DrawImage(texture.BottomLeft, opts)

	// Bottom right corner
	opts.GeoM.Translate(w+3*config.UIScaling, 0)
	screen.DrawImage(texture.BottomRight, opts)

	// Left side
	opts.GeoM.Reset()
	opts.GeoM.Scale(config.UIScaling, h)
	opts.GeoM.Translate(x, y+3*config.UIScaling)
	screen.DrawImage(texture.Left, opts)

	// Right side
	opts.GeoM.Translate(w+3*config.UIScaling, 0)
	screen.DrawImage(texture.Right, opts)

	// Top side
	opts.GeoM.Reset()
	opts.GeoM.Scale(w, config.UIScaling)
	opts.GeoM.Translate(x+3*config.UIScaling, y)
	screen.DrawImage(texture.Top, opts)

	// Bottom side
	opts.GeoM.Translate(0, h+3*config.UIScaling)
	screen.DrawImage(texture.Bottom, opts)

	// Center (background)
	opts.GeoM.Reset()
	opts.GeoM.Scale(w, h)
	opts.GeoM.Translate(x+3*config.UIScaling, y+3*config.UIScaling)
	screen.DrawImage(texture.Center, opts)
}

func DrawTooltipBackground(screen *ebiten.Image, x, y, w, h float64) {
	draw_tooltip(screen, tooltip, x, y, w, h)
}
func DrawButtonBackground(screen *ebiten.Image, hover bool, x, y, w, h float64) {
	if hover {
		draw_tooltip(screen, tooltip_button_hover, x, y, w, h)
	} else {
		draw_tooltip(screen, tooltip_button, x, y, w, h)
	}
}
func DrawInputBackground(screen *ebiten.Image, focused bool, x, y, w, h float64) {
	if focused {
		draw_tooltip(screen, tooltip_input_focused, x, y, w, h)
	} else {
		draw_tooltip(screen, tooltip_input, x, y, w, h)
	}
}
