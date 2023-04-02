package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(GrassBlockState{})
}

type GrassBlockState struct {
	ConnectedBlockState
	CollidableBlockState
}

type GrassBlock struct {
	connectedBlock
	collidableBlock
}

func NewGrassBlock() *GrassBlock {
	return &GrassBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: Grass,
			},
			tex: asset_loader.ConnectedTexture("grass", true, true, true, true),
			connectsTo: []types.BlockType{
				Grass, ShortGrass, TallGrass, Flowers,
				PineTree,
				RedMushroom, WhiteMushroom,
				Stone,
			},
		},
		collidableBlock: collidableBlock{
			collidable:  false,
			playerSpeed: 1,
		},
	}
}

func (b *GrassBlock) State() interface{} {
	return GrassBlockState{
		ConnectedBlockState:  b.connectedBlock.State().(ConnectedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *GrassBlock) LoadState(s interface{}) {
	state := s.(GrassBlockState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
