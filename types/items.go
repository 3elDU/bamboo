package types

import "github.com/hajimehoshi/ebiten/v2"

type ItemType uint

const (
	BlockItem ItemType = iota // Kept for backwards compatibility of IDs
	TestItem                  // Kept for backwards compatibility of IDs
	PineSaplingItem
	StickItem
	FlintItem
	BerryItem
	ClayItem
)

func NewItem(id ItemType) Item {
	switch id {
	case PineSaplingItem:
		return NewPineSaplingItem()
	case StickItem:
		return NewStickItem()
	case FlintItem:
		return NewFlintItem()
	case BerryItem:
		return NewBerryItem()
	case ClayItem:
		return NewClayItem()
	}

	return nil
}

var (
	NewPineSaplingItem func() Item
	NewStickItem       func() Item
	NewFlintItem       func() Item
	NewBerryItem       func() Item
	NewClayItem        func() Item
)

type Item interface {
	Name() string
	Description() string

	Texture() *ebiten.Image
	Type() ItemType
	// Item's hash is calculated from Item's state.
	// As long as two items with the same type have different hashes, they won't stack.
	Hash() uint64

	State() interface{}
	LoadState(interface{})

	Use(pos Vec2u)
}

type BurnableItem interface {
	BurningEnergy() float64
}
