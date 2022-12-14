/*
	Implementations for various block types
*/

package world

import (
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
	Short_Grass
	Tall_Grass
	Flowers
	PineTree
	RedMushroom
	WhiteMushroom
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
		return NewSandBlock(false)
	case Grass:
		return NewGrassBlock()
	case Snow:
		return NewSnowBlock()
	case Short_Grass:
		return NewShortGrassBlock()
	case Tall_Grass:
		return NewTallGrassBlock()
	case Flowers:
		return NewFlowersBlock()
	case PineTree:
		return NewPineTreeBlock()
	case RedMushroom:
		return NewRedMushroomBlock()
	case WhiteMushroom:
		return NewWhiteMushroomBlock()
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
		tex:        asset_loader.ConnectedTexture("grass", true, true, true, true),
		connectsTo: []BlockType{Grass, Short_Grass, Tall_Grass, Flowers, PineTree, RedMushroom, WhiteMushroom, Stone},
	}
}

func NewShortGrassBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1,
			blockType:   Short_Grass,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("short_grass"),
			rotation: 0,
		},
	}
}

func NewTallGrassBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1,
			blockType:   Tall_Grass,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("tall_grass"),
			rotation: 0,
		},
	}
}

func NewFlowersBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1,
			blockType:   Flowers,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture("flowers"),
			rotation: 0,
		},
	}
}

func NewSandBlock(stones bool) *compositeBlock {
	texVariant := "sand"
	if stones {
		texVariant = "sand-stones"
	}

	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.8,
			blockType:   Sand,
		},
		texturedBlock: texturedBlock{
			tex:      asset_loader.Texture(texVariant),
			rotation: float64(util.RandomChoice([]int{0, 90, 180, 270})),
		},
	}
}

func NewWaterBlock() *connectedBlock {
	return &connectedBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 0.2,
			blockType:   Water,
		},
		tex:        asset_loader.ConnectedTexture("lake", false, false, false, false),
		connectsTo: []BlockType{Water},
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
			collidable: true,
			collisionPoints: [4]util.Coords2f{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
			},
			playerSpeed: 0.3,
			blockType:   Stone,
		},
		tex:        asset_loader.ConnectedTexture("stone", false, false, false, false),
		connectsTo: []BlockType{Stone},
	}
}

func NewPineTreeBlock() *connectedBlock {
	return &connectedBlock{
		baseBlock: baseBlock{
			collidable: true,
			blockType:  PineTree,
		},
		tex:        asset_loader.ConnectedTexture("pine", false, false, false, false),
		connectsTo: []BlockType{PineTree},
	}
}

func NewRedMushroomBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1.0,
			blockType:   RedMushroom,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("red-mushroom"),
		},
	}
}

func NewWhiteMushroomBlock() *compositeBlock {
	return &compositeBlock{
		baseBlock: baseBlock{
			collidable:  false,
			playerSpeed: 1.0,
			blockType:   WhiteMushroom,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("white-mushroom"),
		},
	}
}
