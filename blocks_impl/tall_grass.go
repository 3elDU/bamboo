package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/assets"
)

func init() {
	gob.Register(TallGrassState{})
	types.NewTallGrassBlock = NewTallGrassBlock
}

type TallGrassState struct {
	BaseBlockState
	TexturedBlockState
}

type TallGrassBlock struct {
	baseBlock
	texturedBlock
}

func NewTallGrassBlock() types.Block {
	return &TallGrassBlock{
		baseBlock: baseBlock{
			blockType: types.TallGrassBlock,
		},
		texturedBlock: texturedBlock{
			tex:      assets.Texture("tall_grass"),
			rotation: 0,
		},
	}
}

func (b *TallGrassBlock) State() interface{} {
	return TallGrassState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *TallGrassBlock) LoadState(s interface{}) {
	state := s.(TallGrassState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
