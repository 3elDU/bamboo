package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(PineTreeState{})
}

type PineTreeState struct {
	ConnectedBlockState
	CollidableBlockState
}

type pineTree struct {
	connectedBlock
	collidableBlock
}

func NewPineTreeBlock() *pineTree {
	return &pineTree{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: PineTree,
			},
			tex:        asset_loader.ConnectedTexture("pine", false, false, false, false),
			connectsTo: []types.BlockType{PineTree},
		},
		collidableBlock: collidableBlock{
			collidable:      true,
			collisionPoints: defaultCollisionPoints(),
		},
	}
}

func (b pineTree) State() interface{} {
	return PineTreeState{
		ConnectedBlockState:  b.connectedBlock.State().(ConnectedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *pineTree) LoadState(s interface{}) {
	state := s.(PineTreeState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
