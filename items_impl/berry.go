package items_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	types.NewBerryItem = NewBerryItem
	gob.Register(BerryItemState{})
}

type BerryItemState struct {
	BaseItemState
}

type BerryItem struct {
	baseItem
	texture types.Texture
}

func NewBerryItem() types.Item {
	return &BerryItem{
		baseItem: baseItem{
			id: types.BerryItem,
		},
		texture: assets.Texture("cherry"),
	}
}

func (berry *BerryItem) Name() string {
	return "Berry"
}

func (berry *BerryItem) Description() string {
	return "Berry tasty!"
}

func (berry *BerryItem) Texture() *ebiten.Image {
	return berry.texture.Texture()
}

func (berry *BerryItem) State() interface{} {
	return BerryItemState{
		BaseItemState: berry.baseItem.State().(BaseItemState),
	}
}

func (berry *BerryItem) LoadState(s interface{}) {
	state := s.(BerryItemState)
	berry.baseItem.LoadState(state.BaseItemState)
}
