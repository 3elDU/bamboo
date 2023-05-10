package types

import "github.com/hajimehoshi/ebiten/v2"

type ItemType uint

const (
	BlockItem ItemType = iota
	TestItem
)

var NewItem func(id ItemType) Item
var (
	NewBlockItem func(block DrawableBlock) Item
	NewTestItem  func() Item
)

type Item interface {
	Texture() *ebiten.Image
	Type() ItemType
	// Item's hash is calculated from Item's state.
	// As long as two items with the same type have different hashes, they won't stack.
	Hash() uint64

	State() interface{}
	LoadState(interface{})

	Use(world World, pos Vec2u)
}
