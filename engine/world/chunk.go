package world

import (
	"fmt"

	"github.com/3elDU/bamboo/engine/scene_manager"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chunk struct {
	// those are chunk coordinates, not block coordinates
	x, y   int64
	blocks [16][16]Block

	Texture *ebiten.Image

	// Whether a chunk has been modified since last update
	modified bool
	// similar to modified, but indicates that a redraw is required
	// resets on Chunk.Render()
	needsRedraw  bool
	lastAccessed uint64
}

// NewChunk creates new empty Chunk at specified chunk coordinates
func NewChunk(cx, cy int64) *Chunk {
	return &Chunk{
		x: cx, y: cy,
		Texture:      ebiten.NewImage(256, 256),
		modified:     true,
		needsRedraw:  true,
		lastAccessed: scene_manager.Ticks(),
	}
}

// Returns a chunk filled with water
func NewDummyChunk(cx, cy int64) *Chunk {
	c := NewChunk(cx, cy)

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			c.SetBlock(x, y, NewWaterBlock())
		}
	}

	// avoid saving dummy chunk to disk
	c.modified = false

	return c
}

func (c *Chunk) Update(world *World) {
	if c.modified {
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				c.blocks[x][y].Update(world)
			}
		}
	}
}

func (c Chunk) BlockCoords() util.Coords2i {
	return util.Coords2i{X: c.x * 16, Y: c.y * 16}
}

func (c Chunk) Coords() util.Coords2i {
	return util.Coords2i{X: c.x, Y: c.y}
}

func (c *Chunk) At(x, y int) (Block, error) {
	if x < 0 || y < 0 || x > 16 || y > 16 {
		return nil, fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}
	c.lastAccessed = scene_manager.Ticks()
	return c.blocks[x][y], nil
}

func (c *Chunk) SetBlock(x, y int, block Block) error {
	if x > 15 || y > 15 {
		return fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}

	block.SetParentChunk(c)
	block.SetCoords(util.Coords2i{X: c.x*16 + int64(x), Y: c.y*16 + int64(y)})
	c.blocks[x][y] = block
	c.lastAccessed = scene_manager.Ticks()
	c.modified = true
	c.needsRedraw = true
	return nil
}
