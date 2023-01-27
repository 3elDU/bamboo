package world

import (
	"log"
	"math"

	"github.com/3elDU/bamboo/config"
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
			block, err := c.At(x, y)
			if err != nil {
				log.Panicf("Chunk.Render() - chunk.At() failed with %v", err)
			}

			drawableBlock, ok := block.(types.DrawableBlock)
			if !ok {
				continue
			}

			drawableBlock.Render(world, c.Texture(), types.Coords2f{
				X: float64(x) * 16,
				Y: float64(y) * 16,
			})
		}
	}

	c.needsRedraw = false
}

func (world *World) Render(screen *ebiten.Image, playerX, playerY, scaling float64) {
	var (
		screenWidth, screenHeight = screen.Size()
		screenWidthInChunks       = float64(screenWidth) / 256 / scaling
		screenHeightInChunks      = float64(screenHeight) / 256 / scaling
		opts                      = &ebiten.DrawImageOptions{}
	)

	// player is displayed in center of the screen
	// but internally, player coordinates actually represent upper-left corner of the screen
	// so, we need to adjust camera position a bit, so that the camera will be showing the right area
	// hence, we subtract half of screen size, converted to blocks.
	for x := playerX - screenWidthInChunks/2*16 - 16; x < playerX+screenWidthInChunks/2*16+16; x += 16 {
		for y := playerY - screenHeightInChunks/2*16 - 16; y < playerY+screenHeightInChunks/2*16+16; y += 16 {
			// If a chunk is out of world borders, skip it
			if x < 0 || x > float64(config.WorldWidth) || y < 0 || y > float64(config.WorldHeight) {
				continue
			}

			chunk := world.ChunkAtB(uint64(x), uint64(y))
			chunk.Render(world)

			var (
				screenX = (x-playerX-math.Mod(x, 16))*16 + float64(screenWidth)/2 - (float64(screenWidth)/scaling*(scaling-1))/2
				screenY = (y-playerY-math.Mod(y, 16))*16 + float64(screenHeight)/2 - (float64(screenHeight)/scaling*(scaling-1))/2
			)

			opts.GeoM.Reset()
			opts.GeoM.Translate(screenX, screenY)
			opts.GeoM.Scale(scaling, scaling)
			screen.DrawImage(chunk.Texture(), opts)
		}
	}
}
