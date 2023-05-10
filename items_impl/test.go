package items_impl

import (
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

func init() {
	types.NewTestItem = NewTestItem
}

type TestItem struct {
	baseItem
	Tex types.Texture
}

func NewTestItem() types.Item {
	return &TestItem{
		baseItem: baseItem{
			id: types.TestItem,
		},
		Tex: asset_loader.Texture("test_item"),
	}
}

func (item *TestItem) Texture() *ebiten.Image {
	return item.Tex.Texture()
}

func (item *TestItem) Hash() uint64 {
	return 0
}

func (item *TestItem) Use(world types.World, pos types.Vec2u) {
	world.SetBlock(pos.X, pos.Y, types.NewSandBlock(rand.Intn(10) < 5))
}
