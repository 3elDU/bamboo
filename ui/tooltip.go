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

var corner_top_left, corner_top_right, corner_bottom_left, corner_bottom_right *ebiten.Image
var side_left, side_right, side_top, side_bottom *ebiten.Image
var center *ebiten.Image

func init() {
	// Load all the textures

	tooltip := asset_loader.Texture("tooltip").Texture()

	corner_top_left = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 3, Y: 3},
	}))
	corner_top_right = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 4, Y: 0},
		Max: image.Point{X: 7, Y: 3},
	}))
	corner_bottom_left = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 4},
		Max: image.Point{X: 3, Y: 7},
	}))
	corner_bottom_right = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 4, Y: 4},
		Max: image.Point{X: 7, Y: 7},
	}))

	side_left = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 3},
		Max: image.Point{X: 3, Y: 4},
	}))
	side_right = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 4, Y: 3},
		Max: image.Point{X: 7, Y: 4},
	}))
	side_top = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 3, Y: 0},
		Max: image.Point{X: 4, Y: 3},
	}))
	side_bottom = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 3, Y: 4},
		Max: image.Point{X: 4, Y: 7},
	}))

	center = ebiten.NewImageFromImage(tooltip.SubImage(image.Rectangle{
		Min: image.Point{X: 3, Y: 3},
		Max: image.Point{X: 4, Y: 4},
	}))
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
		BottomRight: {X: cursorX, Y: cursorY},
		BottomLeft:  {X: cursorX - int(width) - int(6*config.UIScaling), Y: cursorY},
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

	DrawBackground(screen, x, y, w, h)

	font.RenderFont(screen, text, x+3*config.UIScaling, y+3*config.UIScaling, colors.White)
}

func DrawBackground(screen *ebiten.Image, x, y, w, h float64) {
	opts := &ebiten.DrawImageOptions{}

	// Top left corner
	opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	opts.GeoM.Translate(x, y)
	screen.DrawImage(corner_top_left, opts)

	// Top right corner
	opts.GeoM.Translate(w+3*config.UIScaling, 0)
	screen.DrawImage(corner_top_right, opts)

	// Bottom left corner
	opts.GeoM.Reset()
	opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	opts.GeoM.Translate(x, y+h+3*config.UIScaling)
	screen.DrawImage(corner_bottom_left, opts)

	// Bottom right corner
	opts.GeoM.Translate(w+3*config.UIScaling, 0)
	screen.DrawImage(corner_bottom_right, opts)

	// Left side
	opts.GeoM.Reset()
	opts.GeoM.Scale(config.UIScaling, h)
	opts.GeoM.Translate(x, y+3*config.UIScaling)
	screen.DrawImage(side_left, opts)

	// Right side
	opts.GeoM.Translate(w+3*config.UIScaling, 0)
	screen.DrawImage(side_right, opts)

	// Top side
	opts.GeoM.Reset()
	opts.GeoM.Scale(w, config.UIScaling)
	opts.GeoM.Translate(x+3*config.UIScaling, y)
	screen.DrawImage(side_top, opts)

	// Bottom side
	opts.GeoM.Translate(0, h+3*config.UIScaling)
	screen.DrawImage(side_bottom, opts)

	// Center (background)
	opts.GeoM.Reset()
	opts.GeoM.Scale(w, h)
	opts.GeoM.Translate(x+3*config.UIScaling, y+3*config.UIScaling)
	screen.DrawImage(center, opts)
}
