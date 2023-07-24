/*
	Base structures and interfaces
*/

package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ComponentStyle struct {
	Modified bool

	TextColor  color.Color
	TextShadow bool
	TextSize   float64
}

var defaultStyle = ComponentStyle{
	Modified:   false,
	TextColor:  color.Black,
	TextShadow: true,
	TextSize:   1.0,
}

// Component is the base interface for all components
type Component interface {
	SetParent(parent Component)
	Children() []Component

	// Alignment of the child inside the parent container
	Alignment() ComponentAlignment
	// MaxSize returns maximum space that the component could theoretically occupy
	MaxSize() (float64, float64)
	// Returns actual size of the component
	ComputedSize() (float64, float64)
	// Returns available space that the child could take in its parent, considering the parent's current size.
	CapacityForChild(child Component) (float64, float64)
	// Returns practical maximum space, that the child could take in its parent
	// For most component's this is the same as CapacityForChild
	MaxCapacityForChild(child Component) (float64, float64)

	// Update is pretty much useless,
	// because for many widgets to update, we need to know their size,
	// and their size can't be calculated without access to graphics.
	// So, instead, all those operations are done in Draw() function.
	// But, things not dependent on graphics can be done there
	Update() error
	Draw(screen *ebiten.Image, x, y float64) error

	// ID returns unique identifier of the component, so it can be compared to others
	ID() uint64

	Style() *ComponentStyle
	HasCustomStyles() bool
}

type TextView interface {
	Component

	Text() string
	SetText(text string)
}

// FocusView is an interface for components that can be focused
type FocusView interface {
	Component

	SetFocused(focus bool)
	Focused() bool
}

// ButtonView is an interface for components that can be clicked
type ButtonView interface {
	Component

	IsPressed() bool
	Press()
}

// InputView is an interface for components that accept keyboard input
type InputView interface {
	FocusView

	Input() string
	SetInput(input string)
}

// Used to hint the parent container how the component should be aligned inside it
type ComponentAlignment int

const (
	AlignNone ComponentAlignment = iota
	AlignStart
	AlignCenter
	AlignEnd
)

// baseComponent implements some common methods of View to reduce code repetition
type baseComponent struct {
	parent    Component
	alignment ComponentAlignment
	style     ComponentStyle
	id        uint64
}

var _id uint64 = 0

func newBaseComponent() baseComponent {
	v := baseComponent{
		parent:    nil,
		alignment: AlignNone,
		style:     defaultStyle,
		id:        _id,
	}
	_id += 1
	return v
}
func (b *baseComponent) ID() uint64 {
	return b.id
}
func (b *baseComponent) SetParent(parent Component) {
	b.parent = parent
}
func (b *baseComponent) Align(alignment ComponentAlignment) *baseComponent {
	b.alignment = alignment
	return b
}
func (b *baseComponent) Alignment() ComponentAlignment {
	return b.alignment
}

func (b *baseComponent) Style() *ComponentStyle {
	if !b.HasCustomStyles() {
		return b.parent.Style()
	}
	return &b.style
}
func (b *baseComponent) HasCustomStyles() bool {
	return b.style.Modified
}

// baseFocusView implements some common methods of FocusView to reduce code repetition
type baseFocusView struct {
	focused bool
}

func (b *baseFocusView) SetFocused(focused bool) {
	b.focused = focused
}
func (b *baseFocusView) Focused() bool {
	return b.focused
}
