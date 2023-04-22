package types

import "github.com/hajimehoshi/ebiten/v2"

type ItemType uint

type Item interface {
	Texture() *ebiten.Image
	Type() ItemType

	State() interface{}
	LoadState(interface{})

	Use(world World, pos Vec2u)
}
