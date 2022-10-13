/*
	Various block types, used in the game
*/

package world

import (
	"image/color"
	"math"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Base structure inherited by all blocks
// Contains some basic parameters, so we don't have to implemt them for ourselves
type baseBlock struct {
	// Usually you don't have to set this for youself,
	// Since world.Gen() sets them automatically
	parentChunk *Chunk
	x, y        int

	// Whether collision will work with this block
	collidable bool

	// How fast player could move through this block
	// Calculated by basePlayerSpeed * playerSpeed
	// Applicable only if collidable is false
	playerSpeed float64
}

func (b *baseBlock) Coords() util.Coords2i {
	return util.Coords2i{X: int64(b.x), Y: int64(b.y)}
}

func (b *baseBlock) SetCoords(coords util.Coords2i) {
	if coords.X > 15 || coords.Y > 15 {
		return
	}

}

func (b *baseBlock) ParentChunk() *Chunk {
	return b.parentChunk
}

func (b *baseBlock) SetParentChunk(c *Chunk) {
	b.parentChunk = c
}

func (b *baseBlock) Collidable() bool {
	return b.collidable
}

// Another base structure, to simplify things
type texturedBlock struct {
	tex      *ebiten.Image
	rotation float64 // in degrees
}

func (b *texturedBlock) Render(screen *ebiten.Image, pos util.Coords2f) {
	opts := &ebiten.DrawImageOptions{}

	if b.rotation != 0 {
		w, h := b.tex.Size()
		// Move image half a texture size, so that rotation origin will be in the center
		opts.GeoM.Translate(float64(-w/2), float64(-h/2))
		opts.GeoM.Rotate(b.rotation * (math.Pi / 180))
		pos.X += float64(w / 2)
		pos.Y += float64(h / 2)
	}

	opts.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(b.tex, opts)
}

type coloredBlock struct {
	baseBlock

	color color.Color
}

func NewColoredBlock(color color.Color) *coloredBlock {
	block := coloredBlock{
		baseBlock: baseBlock{},
		color:     color,
	}

	return &block
}

func (b *coloredBlock) Update() {

}

func (b *coloredBlock) Render(screen *ebiten.Image, pos util.Coords2f) {
	// engine.GlobalEngine.FillRectF(float32(target.X), float32(target.Y), 16, 16, b.color)
	ebitenutil.DrawRect(screen, pos.X, pos.Y, 16, 16, b.color)
}

type grassBlock struct {
	*baseBlock
	*texturedBlock
}

func NewGrassBlock() *grassBlock {
	return &grassBlock{
		baseBlock: &baseBlock{
			collidable:  false,
			playerSpeed: 1,
		},
		texturedBlock: &texturedBlock{
			tex:      asset_loader.Texture("grass"),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
	}
}

func (b *grassBlock) Update() {

}

type sandBlock struct {
	*baseBlock
	*texturedBlock
}

func NewSandBlock() *sandBlock {
	return &sandBlock{
		baseBlock: &baseBlock{
			collidable:  false,
			playerSpeed: 0.8,
		},
		texturedBlock: &texturedBlock{
			tex:      asset_loader.Texture("sand"),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
	}
}

func (b *sandBlock) Update() {

}

type waterBlock struct {
	*baseBlock
	*texturedBlock
}

func NewWaterBlock() *waterBlock {
	return &waterBlock{
		baseBlock: &baseBlock{
			collidable:  false,
			playerSpeed: 0.4,
		},
		texturedBlock: &texturedBlock{
			tex: util.RandomChoice([]*ebiten.Image{
				asset_loader.Texture("water1"),
				asset_loader.Texture("water2"),
			}),
			rotation: 0,
		},
	}
}

func (b *waterBlock) Update() {

}

type snowBlock struct {
	*baseBlock
	*texturedBlock
}

func NewSnowBlock() *snowBlock {
	return &snowBlock{
		baseBlock: &baseBlock{
			collidable:  false,
			playerSpeed: 0.7,
		},
		texturedBlock: &texturedBlock{
			tex:      asset_loader.Texture("snow"),
			rotation: 0,
		},
	}
}

func (b *snowBlock) Update() {

}

type stoneBlock struct {
	*baseBlock
	*texturedBlock
}

func NewStoneBlock(height float64) *stoneBlock {
	var texVariant string
	// use different texture depending on mountain height
	switch {
	case height <= 1.51:
		texVariant = "stone1"
	case height <= 1.57:
		texVariant = "stone2"
	case height <= 1.65:
		texVariant = "stone3"
	}

	return &stoneBlock{
		baseBlock: &baseBlock{
			collidable:  false,
			playerSpeed: 0.3,
		},
		texturedBlock: &texturedBlock{
			tex:      asset_loader.Texture(texVariant),
			rotation: 0,
		},
	}
}

func (b *stoneBlock) Update() {

}
