package blocks_impl

import (
	"encoding/gob"
	"fmt"
	"math/rand"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(PineSaplingBlockState{})
	types.NewPineSaplingBlock = NewPineSaplingBlock
}

type PineSaplingBlockState struct {
	BaseBlockState
	TexturedBlockState
	Stage int
}

type PineSaplingBlock struct {
	baseBlock
	texturedBlock
	stage int
}

func NewPineSaplingBlock() types.Block {
	return &PineSaplingBlock{
		baseBlock: baseBlock{
			blockType: types.PineSaplingBlock,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("sapling_block1"),
		},
		stage: 1,
	}
}

func (block *PineSaplingBlock) Update(world types.World) {
	// every 1800 ticks (30 seconds) sapling has 50% chance to grow to next stage
	// each stage will take ~60 seconds, so full tree would grow in ~3.5 minutes
	if scene_manager.Ticks()%180 == 0 && rand.Intn(20) == 0 {
		block.stage++
		if block.stage == 5 {
			world.SetBlock(uint64(block.x), uint64(block.y), types.NewPineTreeBlock())
		} else {
			block.tex = asset_loader.Texture(fmt.Sprintf("sapling_block%v", block.stage))
		}
	}
}

func (block *PineSaplingBlock) Break() {
	types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewPineSaplingItem(),
		Quantity: 1,
	})
	types.GetCurrentWorld().SetBlock(uint64(block.x), uint64(block.y), types.NewGrassBlock())
}

func (block *PineSaplingBlock) State() interface{} {
	return PineSaplingBlockState{
		BaseBlockState:     block.baseBlock.State().(BaseBlockState),
		TexturedBlockState: block.texturedBlock.State().(TexturedBlockState),
		Stage:              block.stage,
	}
}

func (block *PineSaplingBlock) LoadState(s interface{}) {
	state := s.(PineSaplingBlockState)
	block.baseBlock.LoadState(state.BaseBlockState)
	block.texturedBlock.LoadState(state.TexturedBlockState)
	block.stage = state.Stage
}
