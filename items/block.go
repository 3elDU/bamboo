package items

import (
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

/*
	An item, that represents a block, that can be placed. As simple as that
*/

type ItemFromBlockState struct {
	Texture   string
	BlockType types.BlockType
}

type ItemFromBlock struct {
	baseItem
	blockType types.BlockType
	texture   types.Texture
}

func NewItemFromBlock(block types.DrawableBlock) *ItemFromBlock {
	return &ItemFromBlock{
		baseItem: baseItem{
			id: types.ItemType(block.Type()),
		},
		texture:   asset_loader.Texture(block.TextureName()),
		blockType: block.Type(),
	}
}

func (i *ItemFromBlock) Texture() *ebiten.Image {
	return i.texture.Texture()
}

func (i *ItemFromBlock) Use(world types.World, pos types.Vec2u) {
	world.ChunkAtB(pos.X, pos.Y).
		SetBlock(uint(pos.X%16), uint(pos.Y%16), blocks.GetBlockByID(i.blockType))
}

func (i *ItemFromBlock) State() interface{} {
	return ItemFromBlockState{
		Texture:   i.texture.Name(),
		BlockType: i.blockType,
	}
}

func (i *ItemFromBlock) LoadState(s interface{}) {
	state := s.(ItemFromBlockState)
	i.texture = asset_loader.Texture(state.Texture)
	i.blockType = state.BlockType
}
