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

func (b *PineTreeBlock) Break() {
	addedSaplingItem := types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewPineSaplingItem(),
		Quantity: uint8(1 + rand.Intn(2)),
	})

	// add a stick to inventory with 50% chance
	if rand.Float64() > 0.5 {
		types.GetInventory().AddItem(types.ItemSlot{
			Item:     types.NewStickItem(),
			Quantity: 1,
		})
	}

	// if inventory is full, do not do anything
	if !addedSaplingItem {
		return
	}
	types.GetCurrentWorld().SetBlock(uint64(b.x), uint64(b.y), types.NewGrassBlock())
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
