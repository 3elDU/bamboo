package world

import (
	"math"

	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Chunk) Render() {
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			stack, _ := c.At(x, y)

			for _, block := range []Block{stack.bottom, stack.ground, stack.top} {
				block.Render(c.Texture, util.Coords2f{
					X: float64(x) * 16,
					Y: float64(y) * 16,
				})
			}
		}
	}
}

func (world *World) Render(screen *ebiten.Image, px, py, scaling float64) {
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
	for x := px - screenWidthInChunks/2*16 - 16; x < px+screenWidthInChunks/2*16+16; x += 16 {
		for y := py - screenHeightInChunks/2*16 - 16; y < py+screenHeightInChunks/2*16+16; y += 16 {
			chunk, err := world.At(x, y)
			if err != nil {
				panic(err)
			}

			chunk.Render()

			var (
				sx = (x-px-math.Mod(x, 16))*16 + float64(screenWidth)/2 - (float64(screenWidth)/scaling*(scaling-1))/2
				sy = (y-py-math.Mod(y, 16))*16 + float64(screenHeight)/2 - (float64(screenHeight)/scaling*(scaling-1))/2
			)

			opts.GeoM.Reset()
			opts.GeoM.Translate(sx, sy)
			opts.GeoM.Scale(scaling, scaling)
			screen.DrawImage(chunk.Texture, opts)
		}
	}
}
