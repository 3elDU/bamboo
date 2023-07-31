package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/assets"
)

func init() {
	gob.Register(SandState{})
	types.NewSandBlock = NewSandBlock
}

type SandState struct {
	BaseBlockState
}

type SandBlock struct {
	connectedBlock
	collidableBlock
}

func NewSandBlock() types.Block {
	return &SandBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.SandBlock,
			},
			tex: assets.ConnectedTexture("sand", true, true, true, true),
			connectsTo: []types.BlockType{
				types.SandBlock,
			},
		},
		collidableBlock: collidableBlock{
			collidable:  false,
			playerSpeed: 0.8,
		},
	}
}

func (b *SandBlock) State() interface{} {
	return SandState{
		BaseBlockState: b.baseBlock.State().(BaseBlockState),
	}
}

func (b *SandBlock) LoadState(s interface{}) {
	state := s.(SandState)
	b.baseBlock.LoadState(state.BaseBlockState)
}
