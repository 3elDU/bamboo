package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
)

func init() {
	gob.Register(MushroomState{})
}

type MushroomState struct {
	BaseBlockState
	TexturedBlockState
}

type MushroomBlock struct {
	baseBlock
	texturedBlock
}

func NewRedMushroomBlock() *MushroomBlock {
	return &MushroomBlock{
		baseBlock: baseBlock{
			blockType: RedMushroom,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("red-mushroom"),
		},
	}
}

func NewWhiteMushroomBlock() *MushroomBlock {
	return &MushroomBlock{
		baseBlock: baseBlock{
			blockType: WhiteMushroom,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("white-mushroom"),
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
