/*
	Various block types, used in the game
*/

package world

import (
	"github.com/3elDU/bamboo/engine"
	"github.com/3elDU/bamboo/util"
	"github.com/veandco/go-sdl2/sdl"
)

type BaseBlock struct {
	parentChunk *chunk
	x, y        int
}

func (b *BaseBlock) Coords() util.Coords2i {
	return util.Coords2i{X: int64(b.x), Y: int64(b.y)}
}

func (b *BaseBlock) SetCoords(coords util.Coords2i) {
	if coords.X > 15 || coords.Y > 15 {
		return
	}

}

func (b *BaseBlock) ParentChunk() *chunk {
	return b.parentChunk
}

func (b *BaseBlock) SetParentChunk(c *chunk) {
	b.parentChunk = c
}

type coloredBlock struct {
	BaseBlock

	color sdl.Color
}

func NewColoredBlock(color sdl.Color) *coloredBlock {
	block := coloredBlock{
		color: color,
	}

	return &block
}

func (b *coloredBlock) Update() {

}

func (b *coloredBlock) Render(target util.Coords2f) {
	engine.GlobalEngine.FillRectF(float32(target.X), float32(target.Y), 16, 16, b.color)
}
