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
	blocks [16][16]BlockStack

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
			c.SetStack(x, y, BlockStack{
				NewEmptyBlock(),
				NewWaterBlock(),
				NewEmptyBlock(),
			})
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
				c.blocks[x][y].Bottom.Update(world)
				c.blocks[x][y].Ground.Update(world)
				c.blocks[x][y].Top.Update(world)
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

func (c *Chunk) At(x, y int) (*BlockStack, error) {
	if x < 0 || y < 0 || x > 16 || y > 16 {
		return nil, fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}
	c.lastAccessed = scene_manager.Ticks()
	return &c.blocks[x][y], nil
}

func (c *Chunk) SetBlock(x, y int, layer Layer, block Block) error {
	if x > 15 || y > 15 || layer > TopLayer {
		return fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}

	block.SetParentChunk(c)
	block.SetCoords(util.Coords2i{X: c.x*16 + int64(x), Y: c.y*16 + int64(y)})
	block.SetLayer(layer)

	switch layer {
	case BottomLayer:
		c.blocks[x][y].Bottom = block
	case GroundLayer:
		c.blocks[x][y].Ground = block
	case TopLayer:
		c.blocks[x][y].Top = block
	}

	c.lastAccessed = scene_manager.Ticks()
	c.modified = true
	c.needsRedraw = true
	return nil
}

func (c *Chunk) SetBottomBlock(x, y int, block Block) error {
	return c.SetBlock(x, y, BottomLayer, block)
}

func (c *Chunk) SetGroundBlock(x, y int, block Block) error {
	return c.SetBlock(x, y, GroundLayer, block)
}

func (c *Chunk) SetTopBlock(x, y int, block Block) error {
	return c.SetBlock(x, y, TopLayer, block)
}

func (c *Chunk) SetStack(x, y int, stack BlockStack) error {
	if x < 0 || y < 0 || x > 15 || y > 15 {
		return fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}

	blockCoordinates := util.Coords2i{X: c.x*16 + int64(x), Y: c.y*16 + int64(y)}

	stack.Bottom.SetParentChunk(c)
	stack.Bottom.SetCoords(blockCoordinates)
	stack.Ground.SetParentChunk(c)
	stack.Ground.SetCoords(blockCoordinates)
	stack.Top.SetParentChunk(c)
	stack.Top.SetCoords(blockCoordinates)

	c.blocks[x][y] = stack
	c.lastAccessed = scene_manager.Ticks()
	c.modified = true
	c.needsRedraw = true
	return nil
}
