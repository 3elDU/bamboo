package blocks_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	types.NewPitBlock = NewPitBlock
}

type PitBlock struct {
	connectedBlock
	collidableBlock
}

func NewPitBlock() types.Block {
	return &PitBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.PitBlock,
			},
			connectsTo: []types.BlockType{types.PitBlock},
			tex:        assets.ConnectedTexture("pit", true, true, true, true),
		},
		collidableBlock: collidableBlock{
			collidable: true,
		},
	}
}

func (b *PitBlock) State() interface{} {
	return b.baseBlock.State()
}
func (b *PitBlock) LoadState(s interface{}) {
	b.baseBlock.LoadState(s)
}
