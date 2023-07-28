package blocks_impl

import (
	"encoding/gob"
	"fmt"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
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
	pieces   int
	burning  bool
	energy   float64
	burntOut bool
}

func NewCampfireBlock() types.Block {
	return &CampfireBlock{
		baseBlock: baseBlock{
			blockType: types.CampfireBlock,
		},
		pieces:   1,
		burning:  false,
		energy:   1,
		burntOut: false,
	}
}

func (campfire *CampfireBlock) Update(world types.World) {
	if campfire.burning {
		// campfire burns 1 energy per minute
		campfire.energy -= 1.0 / 3600

		if campfire.energy < 0 {
			campfire.burning = false
			campfire.burntOut = true
		}
	}
}

func (campfire *CampfireBlock) Render(_ types.World, screen *ebiten.Image, pos types.Vec2f) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(assets.Texture(campfire.TextureName()).Texture(), opts)
}

func (campfire *CampfireBlock) TextureName() string {
	if campfire.burning {
		return "campfire_burning"
	} else if campfire.burntOut {
		return "campfire_ash"
	} else {
		return fmt.Sprintf("campfire%v", campfire.pieces)
	}
}

func (campfire *CampfireBlock) Break() {
	added := types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewStickItem(),
		Quantity: uint8(campfire.pieces),
	})
	if added {
		types.GetCurrentWorld().SetBlock(uint64(campfire.x), uint64(campfire.y), types.NewGrassBlock())
	}
}

func (campfire *CampfireBlock) AddPiece(item types.BurnableItem) {
	if campfire.pieces < 4 {
		campfire.pieces++
	}

	campfire.energy += item.BurningEnergy()
}

func (campfire *CampfireBlock) LightUp() bool {
	if campfire.burning || campfire.pieces != 4 {
		return false
	}

	campfire.burning = true
	return true
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
}
