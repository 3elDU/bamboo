package blocks_impl

import (
	"encoding/gob"
	"fmt"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(CampfireBlockState{})
	types.NewCampfireBlock = NewCampfireBlock
}

type CampfireBlockState struct {
	BaseBlockState
	Pieces   int
	Burning  bool
	Energy   float64
	BurntOut bool
}

type CampfireBlock struct {
	baseBlock
	texturedBlock
	collidableBlock

	pieces   int
	burning  bool
	energy   float64
	burntOut bool
}

func NewCampfireBlock() types.Block {
	campfire := &CampfireBlock{
		baseBlock: baseBlock{
			blockType: types.CampfireBlock,
		},
		texturedBlock:   texturedBlock{},
		collidableBlock: collidableBlock{collidable: true},

		pieces:   1,
		burning:  false,
		energy:   1,
		burntOut: false,
	}
	campfire.updateTexture()
	return campfire
}

func (campfire *CampfireBlock) updateTexture() {
	if campfire.burning {
		campfire.tex = assets.Texture("campfire_burning")
	} else if campfire.burntOut {
		campfire.tex = assets.Texture("campfire_ash")
	} else {
		campfire.tex = assets.Texture(fmt.Sprintf("campfire%v", campfire.pieces))
	}
}

func (campfire *CampfireBlock) Update(world types.World) {
	if campfire.burning {
		// campfire burns 1 energy per minute
		campfire.energy -= 1.0 / 3600

		if campfire.energy < 0 {
			campfire.burning = false
			campfire.burntOut = true
			campfire.updateTexture()
			campfire.parentChunk.MarkAsModified()
		}
	}
}

func (campfire *CampfireBlock) ToolRequiredToBreak() types.ToolFamily {
	return types.ToolFamilyAxe
}
func (campfire *CampfireBlock) ToolStrengthRequired() types.ToolStrength {
	return types.ToolStrengthBareHand
}
func (campfire *CampfireBlock) Break() {
	added := types.GetPlayerInventory().AddItem(types.ItemSlot{
		Item:     types.NewStickItem(),
		Quantity: uint8(campfire.pieces),
	})
	if added {
		types.GetCurrentWorld().SetBlock(uint64(campfire.x), uint64(campfire.y), types.NewGrassBlock())
	}
}

func (campfire *CampfireBlock) AddPiece(item types.IBurnableItem) bool {
	if campfire.pieces < 4 {
		campfire.pieces++
		campfire.updateTexture()
		campfire.parentChunk.MarkAsModified()
		return true
	}
	if campfire.burning {
		campfire.energy += item.BurningEnergy()
		campfire.parentChunk.MarkAsModified()
		return true
	}
	return false
}

func (campfire *CampfireBlock) LightUp() bool {
	if campfire.burning || campfire.pieces != 4 {
		return false
	}

	campfire.burning = true
	campfire.energy = 1
	campfire.updateTexture()
	campfire.parentChunk.MarkAsModified()
	return true
}

func (campfire *CampfireBlock) ExtinguishCampfire() {
	campfire.burning = false
	campfire.parentChunk.MarkAsModified()
}

func (campfire *CampfireBlock) IsLitUp() bool {
	return campfire.burning
}

func (campfire *CampfireBlock) State() interface{} {
	return CampfireBlockState{
		BaseBlockState: campfire.baseBlock.State().(BaseBlockState),
		Pieces:         campfire.pieces,
		Burning:        campfire.burning,
		Energy:         campfire.energy,
		BurntOut:       campfire.burntOut,
	}
}

func (campfire *CampfireBlock) LoadState(s interface{}) {
	state := s.(CampfireBlockState)
	campfire.baseBlock.LoadState(state.BaseBlockState)
	campfire.pieces = state.Pieces
	campfire.burning = state.Burning
	campfire.energy = state.Energy
	campfire.burntOut = state.BurntOut
	campfire.updateTexture()
}
