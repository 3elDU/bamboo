/*
	Base block type.
	Implements basic methods and fields, so we don't have to rewrite this in every block type
*/

package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(BaseBlockState{})
}

type BaseBlockState struct {
	BlockType types.BlockType
}

// Base structure inherited by all blocks
// Contains some basic parameters, so we don't have to implement them for ourselves
type baseBlock struct {
	// Usually you don't have to set this for youself,
	// Since world.Gen() sets them automatically
	parentChunk types.Chunk
	// Block coordinates in world space
	x, y uint

	// Block types are defined in (blocks.go):13
	// Each block must specify it's type, so that we can actually know what the block it is
	// ( Remember, all blocks are the same interface )
	blockType types.BlockType
}

func (b *baseBlock) Coords() types.Coords2u {
	return types.Coords2u{X: uint64(b.x), Y: uint64(b.y)}
}

func (b *baseBlock) SetCoords(coords types.Coords2u) {
	b.x = uint(coords.X)
	b.y = uint(coords.Y)
}

func (b *baseBlock) ParentChunk() types.Chunk {
	return b.parentChunk
}

func (b *baseBlock) SetParentChunk(c types.Chunk) {
	b.parentChunk = c
}

func (b *baseBlock) Type() types.BlockType {
	return b.blockType
}

func (b *baseBlock) Update(_ types.World) {

}

func (b *baseBlock) State() interface{} {
	return BaseBlockState{
		// Collidable:      b.collidable,
		// CollisionPoints: b.collisionPoints,
		// PlayerSpeed:     b.playerSpeed,
		BlockType: b.blockType,
	}
}

func (b *baseBlock) LoadState(s interface{}) {
	state := s.(BaseBlockState)
	b.blockType = state.BlockType
}
