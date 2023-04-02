/*
	Base structures and interfaces
*/

package ui

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// View is the base interface for all components
type View interface {
	SetParent(parent View)
	Children() []View

	// MaxSize returns maximum space that the component could theoretically occupy
	MaxSize() (float64, float64)
	// ComputedSize returns how much space the component takes in practice
	ComputedSize() (float64, float64)
	// CapacityForChild returns practical maximum space, that the child could take in its parent
	CapacityForChild(child View) (float64, float64)

	// Update is pretty much useless,
	// because for many widgets to update, we need to know their size,
	// and their size can't be calculated without access to graphics.
	// So, instead, all those operations are done in Draw() function.
	// But, things not dependent on graphics can be done there
	Update() error
	Draw(screen *ebiten.Image, x, y float64) error

	// ID returns unique identifier of the component, so it can be compared to others
	ID() uint64
}

// FocusView is an interface for components that can be focused
type FocusView interface {
	View

	SetFocused(focus bool)
	Focused() bool
}

// ButtonView is an interface for components that can be clicked
type ButtonView interface {
	View

	IsPressed() bool
	Press()
}

// InputView is an interface for components that accept keyboard input
type InputView interface {
	FocusView

	Input() string
	SetInput(input string)
}

// baseView implements some common methods of View to reduce repeating code
type baseView struct {
	parent View
	id     uint64
}

func newBaseView() baseView {
	v := baseView{
		parent: nil,
		id:     rand.Uint64(),
	}
	return v
}
func (b *baseView) ID() uint64 {
	return b.id
}
func (b *baseView) SetParent(parent View) {
	b.parent = parent
}

// baseFocusView implements some common methods of FocusView to reduce repeating code
type baseFocusView struct {
	focused bool
}

func (b *baseFocusView) SetFocused(focused bool) {
	b.focused = focused
}
func (b *baseFocusView) Focused() bool {
	return b.focused
}
