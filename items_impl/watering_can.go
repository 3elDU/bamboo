package items_impl

import (
	"encoding/gob"
	"fmt"
	"log"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(WateringCanState{})
	types.NewWateringCanItem = NewWateringCanItem
}

type WateringCanState struct {
	BaseItemState
	WaterAmount int
}

type WateringCanItem struct {
	baseItem
	waterAmount int
}

func NewWateringCanItem() types.Item {
	return &WateringCanItem{
		baseItem: baseItem{
			id: types.WateringCanItem,
		},
		waterAmount: 0,
	}
}

func (item *WateringCanItem) Stackable() bool {
	return false
}
func (item *WateringCanItem) Name() string {
	return "Watering can"
}
func (item *WateringCanItem) Description() string {
	return fmt.Sprintf("Water left: %v", item.waterAmount)
}

func (item *WateringCanItem) Texture() *ebiten.Image {
	if item.waterAmount > 0 {
		return assets.Texture("watering_can_with_water").Texture()
	} else {
		return assets.Texture("watering_can").Texture()
	}
}

func (item *WateringCanItem) Family() types.ToolFamily {
	return types.ToolFamilyNone
}
func (item *WateringCanItem) Strength() types.ToolStrength {
	return types.ToolStrengthClay
}
func (item *WateringCanItem) Use(pos types.Vec2u) {
	switch types.GetCurrentWorld().BlockAt(pos.X, pos.Y).Type() {
	case types.WaterBlock:
		item.waterAmount = 5
	default:
		if block, ok := types.GetCurrentWorld().BlockAt(pos.X, pos.Y).(types.ICropBlock); ok {
			if block.NeedsWatering() && item.waterAmount > 0 {
				block.AddWater()
				item.waterAmount -= 1
			}
		}
	}
}

func (item *WateringCanItem) State() interface{} {
	return WateringCanState{
		BaseItemState: item.baseItem.State().(BaseItemState),
		WaterAmount:   item.waterAmount,
	}
}
func (item *WateringCanItem) LoadState(s interface{}) {
	if state, ok := s.(WateringCanState); ok {
		item.baseItem.LoadState(state.BaseItemState)
		item.waterAmount = state.WaterAmount
	} else {
		log.Println("WateringCanItem: invalid state")
	}
}
