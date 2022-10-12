package world

import (
	"fmt"
	"image/color"
	"math"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (c *Chunk) Render(screen *ebiten.Image, target util.Coords2f) {
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			block, _ := c.At(x, y)
			block.Render(screen, util.Coords2f{
				X: target.X + float64(x)*16,
				Y: target.Y + float64(y)*16,
			})
		}
	}
}

func (world *World) Render(screen *ebiten.Image, playerX, playerY float64) {
	w, h := ebiten.WindowSize()

	for x := 0; x <= w+256; x += 16 * 16 {
		for y := 0; y <= h+256; y += 16 * 16 {
			// calculate Chunk coordinates from screen coordinates
			chunkX := int64((playerX + (float64(x) / 16)) / 16)
			chunkY := int64((playerY + (float64(y) / 16)) / 16)
			chunk, err := world.At(chunkX, chunkY)

			screenX := float64(x) - math.Mod(playerX, 16)*16
			screenY := float64(y) - math.Mod(playerY, 16)*16

			if err != nil {
				ebitenutil.DrawRect(screen, screenX, screenY, 256, 256, color.RGBA{R: 255, G: 0, B: 255, A: 255})
			} else {
				chunk.Render(screen, util.Coords2f{
					X: screenX,
					Y: screenY,
				})
			}

			// write some debug info
			engine.RenderFont(screen, asset_loader.DefaultFont(),
				fmt.Sprintf("coords: %v, %v\nx, y: %v, %v\nerr: %v",
					chunkX, chunkY, x, y, err),
				int(screenX), int(screenY), colors.Black,
			)
		}
	}
}
