package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(StoneState{})
	types.NewStoneBlock = NewStoneBlock
}

type StoneState struct {
	ConnectedBlockState
	CollidableBlockState
}

type StoneBlock struct {
	connectedBlock
	collidableBlock
}

func NewStoneBlock() types.Block {
	return &StoneBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.StoneBlock,
			},
			tex:        assets.ConnectedTexture("stone", false, false, false, false),
			connectsTo: []types.BlockType{types.StoneBlock},
		},
		collidableBlock: collidableBlock{
			collidable:      true,
			collisionPoints: defaultCollisionPoints(),
		},
	}
}

func (b *StoneBlock) State() interface{} {
	return StoneState{
		ConnectedBlockState:  b.connectedBlock.State().(ConnectedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *StoneBlock) LoadState(s interface{}) {
	state := s.(StoneState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
