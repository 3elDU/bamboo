/*
	Implementations for various block types
*/

package world

import (
	"fmt"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Empty BlockType = iota
	Stone
	Water
	Sand
	Grass
	Snow

	// Top-layer blocks
	Sand_Stone
	Grass_Plants_Small
	Grass_Flowers
	Grass_Plants
)

type emptyBlock struct {
	baseBlock
}

func (e *emptyBlock) Update() {

}

func (e *emptyBlock) Render(_ *ebiten.Image, _ util.Coords2f) {

}

func NewEmptyBlock() *emptyBlock {
	return &emptyBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1.0,
			blockType:   Empty,
		},
	}
}

func NewGrassBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1,
			blockType:   Grass,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("grass1"),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
	}
}

func NewGrassVegetationBlock(variant BlockType) *compositeBlock {
	var texture string

	switch variant {
	case Grass_Plants_Small:
		texture = "grass2"
	case Grass_Plants:
		texture = "grass4"
	case Grass_Flowers:
		texture = "grass3"
	default:
		panic(fmt.Sprintf("invalid grass vegetation variant - %v", variant))
	}

	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1.0,
			blockType:   variant,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture(texture),
			rotation: 0,
		},
	}
}

func NewSandBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.8,
			blockType:   Sand,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("sand"),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
	}
}

func NewSandStoneBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.8,
			blockType:   Sand_Stone,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("sand-stones"),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
	}
}

func NewWaterBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.4,
			blockType:   Water,
		},
		texturedBlock: texturedBlock{
			tex: util.RandomChoice([]*ebiten.Image{
				asset_loader.Texture("water1"),
				asset_loader.Texture("water2"),
			}),
			rotation: 0,
		},
	}
}

func NewSnowBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.7,
			blockType:   Snow,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("snow"),
			rotation: 0,
		},
	}
}

func NewStoneBlock(height float64) *compositeBlock {
	var texVariant string
	// use different texture depending on mountain height
	switch {
	case height <= 1.51:
		texVariant = "stone1"
	case height <= 1.57:
		texVariant = "stone2"
	case height <= 1.65:
		texVariant = "stone3"
	}

	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.3,
			blockType:   Stone,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture(texVariant),
			rotation: 0,
		},
	}
}
