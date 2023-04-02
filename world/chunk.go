package world

import (
	"log"

	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chunk struct {
	// those are chunk coordinates, not block coordinates
	x, y   uint64
	blocks [16][16]types.Block

	texture *ebiten.Image

	// Whether a chunk has been modified since last update
	modified bool
	// similar to modified, but indicates that redraw is required
	// resets on Chunk.Render()
	needsRedraw  bool
	lastAccessed uint64
}

// NewChunk creates new empty Chunk at specified chunk coordinates
func NewChunk(cx, cy uint64) *Chunk {
	return &Chunk{
		x: cx, y: cy,
		texture:      ebiten.NewImage(256, 256),
		modified:     true,
		needsRedraw:  true,
		lastAccessed: scene_manager.Ticks(),
	}
}

// Returns a chunk filled with water
func NewDummyChunk(cx, cy uint64) *Chunk {
	c := NewChunk(cx, cy)

	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			c.SetBlock(x, y, blocks.NewWaterBlock())
		}
	}

	// avoid saving dummy chunk to disk
	c.modified = false

	return c
}

func (c *Chunk) Update(world types.World) {
	if c.modified {
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				c.blocks[x][y].Update(world)
			}
		}
	}
}

func (c *Chunk) BlockCoords() types.Coords2u {
	return types.Coords2u{X: c.x * 16, Y: c.y * 16}
}

func (c *Chunk) Coords() types.Coords2u {
	return types.Coords2u{X: c.x, Y: c.y}
}

func (c *Chunk) At(x, y uint) types.Block {
	if x > 16 || y > 16 {
		log.Panicf("invalid coordinates: %v, %v", x, y)
	}
	c.lastAccessed = scene_manager.Ticks()
	return c.blocks[x][y]
}

func (c *Chunk) SetBlock(x, y uint, block types.Block) {
	if x > 15 || y > 15 {
		log.Panicf("invalid coordinates: %v, %v", x, y)
	}

	block.SetParentChunk(c)
	block.SetCoords(types.Coords2u{X: c.x*16 + uint64(x), Y: c.y*16 + uint64(y)})
	c.blocks[x][y] = block
	c.lastAccessed = scene_manager.Ticks()
	c.modified = true
	c.needsRedraw = true
}

func (c *Chunk) TriggerRedraw() {
	c.needsRedraw = true
}

func (c *Chunk) Texture() *ebiten.Image {
	return c.texture
}
