package types

import "github.com/hajimehoshi/ebiten/v2"

type ItemType uint

const (
	BlockItem ItemType = iota
	TestItem           // Kept for backwards compatibility of IDs
	PineSaplingItem
)

func NewItem(id ItemType) Item {
	switch id {
	case BlockItem:
		return NewBlockItem(NewEmptyBlock().(DrawableBlock))
	case PineSaplingItem:
		return NewPineSaplingItem()
	}

	return nil
}

var (
	NewBlockItem       func(block DrawableBlock) Item
	NewPineSaplingItem func() Item
)

type Item interface {
	Texture() *ebiten.Image
	Type() ItemType
	// Item's hash is calculated from Item's state.
	// As long as two items with the same type have different hashes, they won't stack.
	Hash() uint64

	State() interface{}
	LoadState(interface{})

	Use(pos Vec2u)
}
