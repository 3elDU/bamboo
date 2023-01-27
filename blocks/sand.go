package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/util"
)

func init() {
	gob.Register(SandState{})
}

type SandState struct {
	BaseBlockState
	TexturedBlockState
	CollidableBlockState
}

type sand struct {
	baseBlock
	texturedBlock
	collidableBlock
}

func NewSandBlock(stones bool) *sand {
	texVariant := "sand"
	if stones {
		texVariant = "sand-stones"
	}

	return &sand{
		baseBlock: baseBlock{
			blockType: Sand,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture(texVariant),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
		collidableBlock: collidableBlock{
			collidable:  false,
			playerSpeed: 0.8,
		},
	}
}

func (b sand) State() interface{} {
	return SandState{
		BaseBlockState:       b.baseBlock.State().(BaseBlockState),
		TexturedBlockState:   b.texturedBlock.State().(TexturedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *sand) LoadState(s interface{}) {
	state := s.(SandState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
