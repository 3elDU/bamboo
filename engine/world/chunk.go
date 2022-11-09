package world

import (
	"fmt"

	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chunk struct {
	// those are coordinates, not block coordinates
	x, y   int64
	blocks [16][16]BlockStack

	Texture *ebiten.Image

	// Whether a chunk has been modified, and should be written to the disk
	modified bool
}

// NewChunk creates new empty Chunk at specified chunk coordinates
func NewChunk(cx, cy int64) *Chunk {
	return &Chunk{
		x: cx, y: cy,
		Texture:  ebiten.NewImage(256, 256),
		modified: true,
	}
}

func (c Chunk) BlockCoords() util.Coords2i {
	return util.Coords2i{X: c.x * 16, Y: c.y * 16}
}

func (c *Chunk) At(x, y int) (*BlockStack, error) {
	if x > 16 || y > 16 {
		return nil, fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}
	return &c.blocks[x][y], nil
}

func (c *Chunk) SetBlock(x, y int, layer Layer, block Block) error {
	if x > 15 || y > 15 || layer > TopLayer {
		return fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}

	block.SetParentChunk(c)
	block.SetCoords(util.Coords2i{X: int64(x), Y: int64(y)})
	block.SetLayer(layer)

	switch layer {
	case BottomLayer:
		c.blocks[x][y].Bottom = block
	case GroundLayer:
		c.blocks[x][y].Ground = block
	case TopLayer:
		c.blocks[x][y].Top = block
	}

	c.modified = true
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
	if x > 15 || y > 15 {
		return fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}

	for _, block := range [3]Block{stack.Bottom, stack.Ground, stack.Top} {
		block.SetParentChunk(c)
		block.SetCoords(util.Coords2i{X: int64(x), Y: int64(y)})
	}

	c.blocks[x][y] = stack
	c.modified = true
	return nil
}
