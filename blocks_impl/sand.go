package blocks_impl

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/util"
)

func init() {
	gob.Register(SandState{})
	types.NewSandBlock = NewSandBlock
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

func NewSandBlock(stones bool) types.Block {
	texVariant := "sand"
	if stones {
		texVariant = "sand-stones"
	}

	return &SandBlock{
		baseBlock: baseBlock{
			blockType: types.SandBlock,
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

func (sand *SandBlock) Break() {
	if sand.tex.Name() != "sand-stones" {
		return
	}
	added := types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewFlintItem(),
		Quantity: 1,
	})
	if added {
		sand.tex = asset_loader.Texture("sand")
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
