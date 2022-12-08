/*
	Implementations for various block types
*/

package world

import (
	"fmt"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/texture"
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

// Returns an empty interface
func GetBlockByID(id BlockType) Block {
	switch id {
	case Empty:
		return NewEmptyBlock()
	case Stone:
		return NewStoneBlock()
	case Water:
		return NewWaterBlock()
	case Sand:
		return NewSandBlock()
	case Grass:
		return NewGrassBlock()
	case Snow:
		return NewSnowBlock()
	case Sand_Stone:
		return NewSandStoneBlock()
	case Grass_Plants_Small, Grass_Flowers, Grass_Plants:
		return NewGrassVegetationBlock(Grass_Plants_Small)
	}

	return NewEmptyBlock()
}

type emptyBlock struct {
	baseBlock
}

func (e *emptyBlock) Update(_ *World) {

}
func (e *emptyBlock) Render(_ *World, _ *ebiten.Image, _ util.Coords2f) {

}
func (e *emptyBlock) TextureName() string {
	return ""
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

func NewGrassBlock() *connectedBlock {
	return &connectedBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1,
			blockType:   Grass,
		},
		tex: asset_loader.ConnectedTexture("grass", true, true, true, true),
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
			tex: util.RandomChoice([]texture.Texture{
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

func NewStoneBlock() *connectedBlock {
	return &connectedBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.3,
			blockType:   Stone,
		},
		tex: asset_loader.ConnectedTexture("stone", false, false, false, false),
	}
}
