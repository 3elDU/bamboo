package items_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewStickItem = NewStickItem
	gob.Register(StickItemState{})
}

type StickItemState struct {
	BaseItemState
	Tex string
}

type StickItem struct {
	baseItem
	Tex types.Texture
}

func NewStickItem() types.Item {
	return &StickItem{
		baseItem: baseItem{
			id: types.StickItem,
		},
		Tex: assets.Texture("stick_item"),
	}
}

func (item *StickItem) Name() string {
	return "Stick"
}

func (item *StickItem) Description() string {
	return ""
}

func (item *StickItem) BurningEnergy() float64 {
	return 1
}

func (item *StickItem) Texture() *ebiten.Image {
	return item.Tex.Texture()
}

func (item *StickItem) Hash() uint64 {
	return uint64(item.id)
}

func (item *StickItem) Use(pos types.Vec2u) {
	b, ok := types.GetCurrentWorld().BlockAt(pos.X, pos.Y).(types.ICampfireBlock)
	if ok {
		// if used on campfire block, add a piece to it
		b.AddPiece(item)
		types.GetInventory().RemoveItem(types.ItemSlot{Item: item, Quantity: 1})
		return
	}

	// campfire can only be placed on empty grass block
	if types.GetCurrentWorld().BlockAt(pos.X, pos.Y).Type() != types.GrassBlock {
		return
	}

	// campfire can't be placed near sand or water
	if types.GetCurrentWorld().BlockNeighboringWith(pos.X, pos.Y, []types.BlockType{types.WaterBlock, types.SandBlock}) {
		return
	}

	types.GetCurrentWorld().SetBlock(pos.X, pos.Y, types.NewCampfireBlock())
	types.GetInventory().RemoveItem(types.ItemSlot{Item: item, Quantity: 1})
}

func (item *StickItem) State() interface{} {
	return StickItemState{
		BaseItemState: item.baseItem.State().(BaseItemState),
		Tex:           item.Tex.Name(),
	}
}

func (item *StickItem) LoadState(s interface{}) {
	state := s.(StickItemState)
	item.baseItem.LoadState(state.BaseItemState)
	item.Tex = assets.Texture(state.Tex)
}
