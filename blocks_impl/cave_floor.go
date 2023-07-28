package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(CaveFloorState{})
	types.NewCaveFloorBlock = NewCaveFloorBlock
}

type CaveFloorState struct {
	BaseBlockState
	TexturedBlockState
}

type CaveFloorBlock struct {
	baseBlock
	texturedBlock
}

func NewCaveFloorBlock() types.Block {
	return &CaveFloorBlock{
		baseBlock: baseBlock{
			blockType: types.CaveFloorBlock,
		},
		texturedBlock: texturedBlock{
			tex:      assets.Texture("cave_floor"),
			rotation: 0,
		},
	}
}

func (b *CaveFloorBlock) State() interface{} {
	return CaveFloorState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *CaveFloorBlock) LoadState(s interface{}) {
	state := s.(CaveFloorState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
