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
	ToolStrengthBareHand ToolStrength = iota // 1
	ToolStrengthWood                         // 4
	ToolStrengthClay                         // 8
	ToolStrengthGold                         // 16
	ToolStrengthCopper                       // 32
	ToolStrengthIron                         // 128
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
	RawIronItem
	IronIngotItem
	ClayPickaxeItem
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
	case RawIronItem:
		return NewRawIronItem()
	case IronIngotItem:
		return NewIronIngotItem()
	case ClayPickaxeItem:
		return NewClayPickaxeItem()
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
	NewRawIronItem     func() Item
	NewIronIngotItem   func() Item
	NewClayPickaxeItem func() Item
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

// An interactive item, such as pickaxe, sword, or any other tool
type IToolItem interface {
	ToolFamily() ToolFamily
	ToolStrength() ToolStrength
	UseTool(pos Vec2u)
}

// An item that produces energy by burning
type IBurnableItem interface {
	BurningEnergy() float64
}

// An item that can be smelted in a furnace
type ISmeltableItem interface {
	SmeltingEnergyRequired() float64
	// This function will be called when the smelting process has ended.
	// It should return the result of the smelting process
	Smelt() Item
}
