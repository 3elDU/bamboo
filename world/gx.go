package world

import (
	"math"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Chunk) Render(world types.World) {
	// do not redraw a chunk, when there is no need to
	if !c.needsRedraw {
		return
	}

	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			block := c.At(x, y)

			drawableBlock, ok := block.(types.DrawableBlock)
			if !ok {
				continue
			}

			drawableBlock.Render(world, c.Texture(), types.Vec2f{
				X: float64(x) * 16,
				Y: float64(y) * 16,
			}, c.recursiveRedraw)
		}
	}

	c.needsRedraw = false
	c.recursiveRedraw = false
}

func BlockToScreen(screen *ebiten.Image, player types.Vec2f, block types.Vec2u, scaling float64) types.Vec2f {
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	return types.Vec2f{
		X: (float64(block.X)-player.X)*16 + float64(screenWidth)/2 - (float64(screenWidth)/scaling*(scaling-1))/2,
		Y: (float64(block.Y)-player.Y)*16 + float64(screenHeight)/2 - (float64(screenHeight)/scaling*(scaling-1))/2,
	}
}

func (world *World) Render(screen *ebiten.Image, playerX, playerY, scaling float64) {
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	screenWidthInChunks := float64(screenWidth) / 256 / scaling
	screenHeightInChunks := float64(screenHeight) / 256 / scaling
	opts := &ebiten.DrawImageOptions{}

	// Adjust camera position to show the right area
	cameraOffsetX := screenWidthInChunks / 2 * 16
	cameraOffsetY := screenHeightInChunks / 2 * 16

	for x := playerX - cameraOffsetX - 16; x < playerX+cameraOffsetX+16; x += 16 {
		for y := playerY - cameraOffsetY - 16; y < playerY+cameraOffsetY+16; y += 16 {
			// Skip chunks that are out of world borders
			if x < 0 || x > float64(world.metadata.Size.X) || y < 0 || y > float64(world.metadata.Size.Y) {
				continue
			}

			chunk := world.ChunkAtB(uint64(x), uint64(y))
			needsRedraw := chunk.(*Chunk).needsRedraw
			chunk.Render(world)

			screenX := (x - playerX - math.Mod(x, 16)) * 16
			screenX += float64(screenWidth)/2 - (float64(screenWidth)/scaling*(scaling-1))/2
			screenY := (y - playerY - math.Mod(y, 16)) * 16
			screenY += float64(screenHeight)/2 - (float64(screenHeight)/scaling*(scaling-1))/2

			opts.GeoM.Reset()
			opts.GeoM.Translate(screenX, screenY)
			opts.GeoM.Scale(scaling, scaling)
			screen.DrawImage(chunk.Texture(), opts)

			if config.DebugMode && needsRedraw {
				font.RenderFontWithOptions(screen, "REDRAW", screenX*scaling, screenY*scaling, colors.C("red"), 1.0, false)
			}
		}
	}
}
