package items_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

func init() {
	gob.Register(PineSaplingItemState{})
	types.NewPineSaplingItem = NewPineSaplingItem
}

type PineSaplingItemState struct {
	BaseItemState
	Tex string
}

type PineSaplingItem struct {
	baseItem
	Tex types.Texture
}

func NewPineSaplingItem() types.Item {
	return &PineSaplingItem{
		baseItem: baseItem{
			id: types.PineSaplingItem,
		},
		Tex: assets.Texture("sapling_item"),
	}
}

func (item *PineSaplingItem) Name() string {
	return "Pine sapling"
}

func (item *PineSaplingItem) Description() string {
	return "Try to plant this in the ground with F"
}

func (item *PineSaplingItem) BurningEnergy() float64 {
	return 0.5
}

func (item *PineSaplingItem) Texture() *ebiten.Image {
	return item.Tex.Texture()
}

func (item *PineSaplingItem) Hash() uint64 {
	return uint64(item.id)
}

func (item *PineSaplingItem) Use(pos types.Vec2u) {
	// sapling can only be planted on grass and it's derivatives
	if !slices.Contains([]types.BlockType{types.GrassBlock, types.ShortGrassBlock, types.FlowersBlock}, types.GetCurrentWorld().BlockAt(pos.X, pos.Y).Type()) {
		return
	}

	// sapling cannot grow near sand
	if types.GetCurrentWorld().BlockNeighboringWith(pos.X, pos.Y, []types.BlockType{types.SandBlock}) {
		return
	}

	types.GetCurrentWorld().SetBlock(pos.X, pos.Y, types.NewPineSaplingBlock())
	types.GetInventory().RemoveItem(types.ItemSlot{
		Item:     item,
		Quantity: 1,
	})
}

func (item *PineSaplingItem) State() interface{} {
	return PineSaplingItemState{
		BaseItemState: item.baseItem.State().(BaseItemState),
		Tex:           item.Tex.Name(),
	}
}

func (item *PineSaplingItem) LoadState(s interface{}) {
	state := s.(PineSaplingItemState)
	item.baseItem.LoadState(state.BaseItemState)
	item.Tex = assets.Texture(state.Tex)
}
