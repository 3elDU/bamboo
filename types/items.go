package types

import "github.com/hajimehoshi/ebiten/v2"

// The tool family that the item belongs to
type ToolFamily int

const (
	ToolFamilyNone ToolFamily = iota
	ToolFamilySword
	ToolFamilyPickaxe
	ToolFamilyAxe
	ToolFamilyShovel
	ToolFamilyScissors
)

// Represents "Hardness" of a material
type ToolStrength int

const (
	ToolStrengthBareHand ToolStrength = iota
	ToolStrengthGold
	ToolStrengthTin
	ToolStrengthWood
	ToolStrengthClay
	ToolStrengthCopper
	ToolStrengthBronze
	ToolStrengthIron
	ToolStrengthAluminium
	ToolStrengthSteel
)

type ItemType uint

const (
	BlockItem ItemType = iota // Kept for backwards compatibility of IDs
	TestItem                  // Kept for backwards compatibility of IDs
	PineSaplingItem
	StickItem
	FlintItem
	BerryItem
	ClayItem
	WateringCanItem
	ClayShovelItem
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
	case WateringCanItem:
		return NewWateringCanItem()
	case ClayShovelItem:
		return NewClayShovelItem()
	}

	return nil
}

var (
	NewPineSaplingItem func() Item
	NewStickItem       func() Item
	NewFlintItem       func() Item
	NewBerryItem       func() Item
	NewClayItem        func() Item
	NewWateringCanItem func() Item
	NewClayShovelItem  func() Item
)

type Item interface {
	Name() string
	Description() string

	Texture() *ebiten.Image
	Type() ItemType
	Stackable() bool

	State() interface{}
	LoadState(interface{})
}

type Tool interface {
	Family() ToolFamily
	Strength() ToolStrength
	Use(pos Vec2u)
}

type IBurnableItem interface {
	BurningEnergy() float64
}
