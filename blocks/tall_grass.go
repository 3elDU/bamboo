package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(TallGrassState{})
}

type TallGrassState struct {
	BaseBlockState
	TexturedBlockState
}

type TallGrassBlock struct {
	baseBlock
	texturedBlock
}

func NewTallGrassBlock() *TallGrassBlock {
	return &TallGrassBlock{
		baseBlock: baseBlock{
			blockType: TallGrass,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("tall_grass"),
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
