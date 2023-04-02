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

type SandBlock struct {
	baseBlock
	texturedBlock
	collidableBlock
}

func NewSandBlock(stones bool) *SandBlock {
	texVariant := "sand"
	if stones {
		texVariant = "sand-stones"
	}

	return &SandBlock{
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

func (b *SandBlock) State() interface{} {
	return SandState{
		BaseBlockState:       b.baseBlock.State().(BaseBlockState),
		TexturedBlockState:   b.texturedBlock.State().(TexturedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *SandBlock) LoadState(s interface{}) {
	state := s.(SandState)
	b.baseBlock.LoadState(state.BaseBlockState)
	b.texturedBlock.LoadState(state.TexturedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
