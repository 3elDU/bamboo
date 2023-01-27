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

type mushroom struct {
	baseBlock
	texturedBlock
}

func NewRedMushroomBlock() *mushroom {
	return &mushroom{
		baseBlock: baseBlock{
			blockType: RedMushroom,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("red-mushroom"),
		},
	}
}

func NewWhiteMushroomBlock() *mushroom {
	return &mushroom{
		baseBlock: baseBlock{
			blockType: WhiteMushroom,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("white-mushroom"),
		},
	}
}

func (b mushroom) State() interface{} {
	return MushroomState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *mushroom) LoadState(s interface{}) {
	state := s.(MushroomState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
}
