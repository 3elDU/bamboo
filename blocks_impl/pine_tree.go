package blocks_impl

import (
	"encoding/gob"
	"math/rand"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(PineTreeState{})
	types.NewPineTreeBlock = NewPineTreeBlock
}

type PineTreeState struct {
	ConnectedBlockState
	CollidableBlockState
}

type PineTreeBlock struct {
	connectedBlock
	collidableBlock
}

func NewPineTreeBlock() types.Block {
	return &PineTreeBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.PineTreeBlock,
			},
			tex:        assets.ConnectedTexture("pine", false, false, false, false),
			connectsTo: []types.BlockType{types.PineTreeBlock},
		},
		collidableBlock: collidableBlock{
			collidable:      true,
			collisionPoints: defaultCollisionPoints(),
		},
	}
}
func (block *PineTreeBlock) ToolRequiredToBreak() types.ToolFamily {
	return types.ToolFamilyAxe
}
func (blok *PineTreeBlock) ToolStrengthRequired() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (b *PineTreeBlock) Break() {
	if types.GetPlayerInventory().AddItems(
		// 1-2 saplings
		types.NewItemSlot(types.NewPineSaplingItem(), uint8(1+rand.Intn(2))),
		// 1-2 sticks
		types.NewItemSlot(types.NewStickItem(), uint8(1+rand.Intn(2))),
	) {
		types.GetCurrentWorld().SetBlock(uint64(b.x), uint64(b.y), types.NewGrassBlock())
	}
}

func (b *PineTreeBlock) State() interface{} {
	return PineTreeState{
		ConnectedBlockState:  b.connectedBlock.State().(ConnectedBlockState),
		CollidableBlockState: b.collidableBlock.State().(CollidableBlockState),
	}
}

func (b *PineTreeBlock) LoadState(s interface{}) {
	state := s.(PineTreeState)
	b.connectedBlock.LoadState(state.ConnectedBlockState)
	b.collidableBlock.LoadState(state.CollidableBlockState)
}
