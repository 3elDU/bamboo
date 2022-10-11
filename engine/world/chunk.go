package world

import (
	"fmt"

	"github.com/3elDU/bamboo/util"
)

type chunk struct {
	// chunk coordinates, not block coordinates
	x, y   int64
	blocks [16][16]Block
}

// Creates new empty chunk with specified coordinates
func NewChunk(x, y int64) *chunk {
	return &chunk{x: x, y: y}
}

func (c *chunk) BlockCoords() util.Coords2i {
	return util.Coords2i{X: c.x * 16, Y: c.y * 16}
}

func (c *chunk) At(x, y int) (Block, error) {
	if x > 16 || y > 16 {
		return nil, fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}
	return c.blocks[x][y], nil
}

func (c *chunk) SetBlock(x, y int, block Block) error {
	if x > 15 || y > 15 {
		return fmt.Errorf("invalid coordinates: %v, %v", x, y)
	}
	block.SetParentChunk(c)
	block.SetCoords(util.Coords2i{X: int64(x), Y: int64(y)})
	c.blocks[x][y] = block
	return nil
}
