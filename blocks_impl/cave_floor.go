package blocks_impl

import (
	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	types.NewCaveFloorBlock = NewCaveFloorBlock
}

type CaveFloorBlock struct {
	baseBlock
	texturedBlock
	hasGrass bool
}

func (b *CaveFloorBlock) updateTexture() {
	if b.hasGrass {
		b.tex = assets.Texture("cave_floor_grass")
	} else {
		b.tex = assets.Texture("cave_floor")
	}
}

func NewCaveFloorBlock(hasGrass bool) types.Block {
	block := &CaveFloorBlock{
		baseBlock: baseBlock{
			blockType: types.CaveFloorBlock,
		},
		texturedBlock: texturedBlock{
			tex:      assets.Texture("cave_floor"),
			rotation: 0,
		},
		hasGrass: hasGrass,
	}
	block.updateTexture()
	return block
}

func (b *CaveFloorBlock) State() interface{} {
	return b.hasGrass
}

func (b *CaveFloorBlock) LoadState(s interface{}) {
	hasGrass := s.(bool)
	b.hasGrass = hasGrass
	b.updateTexture()
}
