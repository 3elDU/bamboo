package items_impl

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"hash/fnv"
)

func init() {
	gob.Register(BlockItemState{})
	types.NewBlockItem = NewBlockItem
}

type BlockItemState struct {
	BlockState interface{}
	BlockType  types.BlockType
}

type BlockItem struct {
	baseItem
	block     types.DrawableBlock
	blockType types.BlockType
	texture   types.Texture
}

func NewBlockItem(block types.DrawableBlock) types.Item {
	return &BlockItem{
		baseItem: baseItem{
			id: types.ItemType(block.Type()),
		},
		block:     block,
		texture:   asset_loader.Texture(block.TextureName()),
		blockType: block.Type(),
	}
}

func (i *BlockItem) Hash() uint64 {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, i.blockType)

	hasher := fnv.New64a()
	hasher.Write(buf.Bytes())
	return hasher.Sum64()
}

func (i *BlockItem) Texture() *ebiten.Image {
	return i.texture.Texture()
}

func (i *BlockItem) Use(world types.World, pos types.Vec2u) {
	// keep the block state, but create a new block instance
	state := i.block.State()
	block := types.NewBlock(i.blockType)
	block.LoadState(state)

	world.ChunkAtB(pos.X, pos.Y).
		SetBlock(uint(pos.X%16), uint(pos.Y%16), block)
}

func (i *BlockItem) State() interface{} {
	return BlockItemState{
		BlockState: i.block.State(),
		BlockType:  i.blockType,
	}
}

func (i *BlockItem) LoadState(s interface{}) {
	state := s.(BlockItemState)
	i.blockType = state.BlockType
	i.block = types.NewBlock(state.BlockType).(types.DrawableBlock)
	i.block.LoadState(state.BlockState)
	i.texture = asset_loader.Texture(i.block.TextureName())
}
