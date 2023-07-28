package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/assets"
)

func init() {
	gob.Register(MushroomState{})
	types.NewRedMushroomBlock = NewRedMushroomBlock
	types.NewWhiteMushroomBlock = NewWhiteMushroomBlock
}

type MushroomState struct {
	BaseBlockState
	TexturedBlockState
}

type MushroomBlock struct {
	baseBlock
	texturedBlock
}

func NewRedMushroomBlock() types.Block {
	return &MushroomBlock{
		baseBlock: baseBlock{
			blockType: types.RedMushroomBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture("red-mushroom"),
		},
	}
}

func NewWhiteMushroomBlock() types.Block {
	return &MushroomBlock{
		baseBlock: baseBlock{
			blockType: types.WhiteMushroomBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture("white-mushroom"),
		},
	}
}

func (b *MushroomBlock) State() interface{} {
	return MushroomState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *MushroomBlock) LoadState(s interface{}) {
	state := s.(MushroomState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
