package types

import "github.com/hajimehoshi/ebiten/v2"

type ItemType uint

const (
	BlockItem ItemType = iota
)

var NewItem func(id ItemType) Item
var (
	NewBlockItem func(block DrawableBlock) Item
)

type Item interface {
	Texture() *ebiten.Image
	Type() ItemType

	State() interface{}
	LoadState(interface{})

	Use(world World, pos Vec2u)
}
