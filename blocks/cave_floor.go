package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(CaveFloorState{})
}

type CaveFloorState struct {
	BaseBlockState
	TexturedBlockState
}

type CaveFloorBlock struct {
	baseBlock
	texturedBlock
}

func NewCaveFloorBlock() *CaveFloorBlock {
	return &CaveFloorBlock{
		baseBlock: baseBlock{
			blockType: CaveFloor,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("cave_floor"),
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
