package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(ShortGrassState{})
}

type ShortGrassState struct {
	BaseBlockState
	TexturedBlockState
}

type ShortGrassBlock struct {
	baseBlock
	texturedBlock
}

func NewShortGrassBlock() *ShortGrassBlock {
	return &ShortGrassBlock{
		baseBlock: baseBlock{
			blockType: ShortGrass,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("short_grass"),
			rotation: 0,
		},
	}
}

func (b *ShortGrassBlock) State() interface{} {
	return ShortGrassState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *ShortGrassBlock) LoadState(s interface{}) {
	state := s.(ShortGrassState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
