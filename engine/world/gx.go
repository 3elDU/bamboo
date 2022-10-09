package world

import (
	"fmt"
	"math"

	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/util"
)

func (c *chunk) Render(target util.Coords2f) {
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			block, _ := c.At(x, y)
			block.Render(util.Coords2f{
				X: target.X + float64(x)*16,
				Y: target.Y + float64(y)*16,
			})
		}
	}
}

func (world *World) Render(playerX, playerY float64) {
	w, h := engine.GlobalEngine.Win.GetSize()

	for x := 0; x <= int(w)+256; x += 16 * 16 {
		for y := 0; y <= int(h)+256; y += 16 * 16 {
			// calculate chunk coordinates from screen coordinates
			chunkX := int64((playerX + (float64(x) / 16)) / 16)
			chunkY := int64((playerY + (float64(y) / 16)) / 16)
			chunk := world.At(chunkX, chunkY)

			var screenX float64 = float64(x) - math.Mod(playerX, 16)*16
			var screenY float64 = float64(y) - math.Mod(playerY, 16)*16

			chunk.Render(util.Coords2f{
				X: screenX,
				Y: screenY,
			})

			// write some debug info
			engine.GlobalEngine.RenderFont(
				asset_loader.Assets.DefaultFont(), int32(screenX), int32(screenY),
				fmt.Sprintf("coords: %v, %v\nx, y: %v, %v",
					chunkX, chunkY, x, y),
				colors.Red,
			)
		}
	}
}
