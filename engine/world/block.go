/*
	Declrations of basic block types
*/

package world

import (
	"fmt"
	"math"

	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Layer int

const (
	BottomLayer Layer = iota
	GroundLayer
	TopLayer
)

type BlockType int

type Block interface {
	Coords() util.Coords2i
	SetCoords(coords util.Coords2i)
	ParentChunk() *Chunk
	SetParentChunk(chunk *Chunk)
	SetLayer(layer Layer)
	Layer() Layer
	Type() BlockType

	// Whether the player should collide with the block
	Collidable() bool

	Update()
	Render(screen *ebiten.Image, pos util.Coords2f)
}

// The map is actually 3D, and consists of three layers:
//
//  1. Fossils / Ore
//  2. Ground block ( the one you`ll see the most )
//  3. Top block - decoration / vegetation / player buildings / etc.
type BlockStack struct {
	bottom Block
	ground Block
	top    Block
}

// Base structure inherited by all blocks
// Contains some basic parameters, so we don't have to implement them for ourselves
type baseBlock struct {
	// Usually you don't have to set this for youself,
	// Since world.Gen() sets them automatically
	parentChunk *Chunk
	x, y        int
	layer       Layer

	// Whether collision will work with this block
	collidable bool

	// How fast player could move through this block
	// Calculated by basePlayerSpeed * playerSpeed
	// Applicable only if collidable is false
	playerSpeed float64

	// Block types are defined in (blocks.go):13
	// Each block must specify it's type, so that we can actually know what the block it is
	// ( Remember, all blocks are the same interface )
	blockType BlockType
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

func (b *baseBlock) Layer() Layer {
	return b.layer
}

func (b *baseBlock) SetLayer(layer Layer) {
	b.layer = layer
}

func (b *baseBlock) Collidable() bool {
	return b.collidable
}

func (b *baseBlock) Type() BlockType {
	return b.blockType
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

// simple block that has derives from both baseBlock and texturedBlock
// used by almost all blocks out there, that don't require any complex behaviour
type compositeBlock struct {
	baseBlock
	texturedBlock
}

func (b *compositeBlock) Update() {
	fmt.Println("hello")
}
