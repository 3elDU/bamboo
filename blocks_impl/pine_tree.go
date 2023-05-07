package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(PineTreeState{})
	types.NewPineTreeBlock = NewPineTreeBlock
}

type PineTreeState struct {
	ConnectedBlockState
	CollidableBlockState
}

type PineTreeBlock struct {
	connectedBlock
	collidableBlock
}

func NewPineTreeBlock() types.Block {
	return &PineTreeBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.PineTreeBlock,
			},
			tex:        asset_loader.ConnectedTexture("pine", false, false, false, false),
			connectsTo: []types.BlockType{types.PineTreeBlock},
		},
		collidableBlock: collidableBlock{
			collidable:      true,
			collisionPoints: defaultCollisionPoints(),
		},
	}
}

func (b *PineTreeBlock) State() interface{} {
	return PineTreeState{
		ConnectedBlockState:  b.connectedBlock.State().(ConnectedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *PineTreeBlock) LoadState(s interface{}) {
	state := s.(PineTreeState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
