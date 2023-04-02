package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(StoneState{})
}

type StoneState struct {
	ConnectedBlockState
	CollidableBlockState
}

type StoneBlock struct {
	connectedBlock
	collidableBlock
}

func NewStoneBlock() *StoneBlock {
	return &StoneBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: Stone,
			},
			tex:        asset_loader.ConnectedTexture("stone", false, false, false, false),
			connectsTo: []types.BlockType{Stone},
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
