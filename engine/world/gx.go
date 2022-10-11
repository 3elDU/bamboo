package world

import (
	"fmt"
	"math"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *chunk) Render(screen *ebiten.Image, target util.Coords2f) {
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
			// calculate chunk coordinates from screen coordinates
			chunkX := int64((playerX + (float64(x) / 16)) / 16)
			chunkY := int64((playerY + (float64(y) / 16)) / 16)
			chunk := world.At(chunkX, chunkY)

			var screenX float64 = float64(x) - math.Mod(playerX, 16)*16
			var screenY float64 = float64(y) - math.Mod(playerY, 16)*16

			chunk.Render(screen, util.Coords2f{
				X: screenX,
				Y: screenY,
			})

			// write some debug info
			engine.RenderFont(screen, asset_loader.GlobalAssets.DefaultFont(),
				fmt.Sprintf("coords: %v, %v\nx, y: %v, %v",
					chunkX, chunkY, x, y),
				int(screenX), int(screenY), colors.Red,
			)
		}
	}
}
