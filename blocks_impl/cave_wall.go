package blocks_impl

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(CaveWallState{})
	types.NewCaveWallBlock = NewCaveWallBlock
}

type CaveWallState struct {
	ConnectedBlockState
	CollidableBlockState
}

type CaveWallBlock struct {
	connectedBlock
	collidableBlock
}

func NewCaveWallBlock() types.Block {
	return &CaveWallBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.CaveWallBlock,
			},
			tex:        asset_loader.ConnectedTexture("cave_wall", false, false, false, false),
			connectsTo: []types.BlockType{types.CaveWallBlock},
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
