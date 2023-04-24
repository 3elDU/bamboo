package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(CaveWallState{})
}

type CaveWallState struct {
	ConnectedBlockState
	CollidableBlockState
}

type CaveWallBlock struct {
	connectedBlock
	collidableBlock
}

func NewCaveWallBlock() *CaveWallBlock {
	return &CaveWallBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: CaveWall,
			},
			tex:        asset_loader.ConnectedTexture("cave_wall", false, false, false, false),
			connectsTo: []types.BlockType{CaveWall},
		},
		collidableBlock: collidableBlock{
			collidable:      true,
			collisionPoints: defaultCollisionPoints(),
		},
	}
}

func (b *CaveWallBlock) State() interface{} {
	return CaveWallState{
		ConnectedBlockState:  b.connectedBlock.State().(ConnectedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *CaveWallBlock) LoadState(s interface{}) {
	state := s.(CaveWallState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
