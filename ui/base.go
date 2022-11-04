/*
	Base structures and interfaces
*/

package ui

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// base interface for all components
type View interface {
	SetParent(parent View)
	Children() []View

	// maximum space that the component could theoretically occupy in their container
	MaxSize() (float64, float64)
	// how much space the component takes in practice
	ComputedSize() (float64, float64)
	// practical maximum space, that provided child could take in the container
	CapacityForChild(child View) (float64, float64)

	// Update is pretty much useless
	// Because for many widgets to update, we need to know their size
	// And their size can't be calculated without access to graphics
	// So, instead, all those operations are done in Draw() function
	// But, things not dependent on graphics can be done there
	Update() error
	Draw(screen *ebiten.Image, x, y float64) error

	// returns unique identifier of the component, so it can be compared to others
	ID() uint64
}

// interface for various elements, that do have a focus
type FocusView interface {
	View

	SetFocused(focus bool)
	Focused() bool
}

type ButtonView interface {
	View

	IsPressed() bool
	Press() // Simulates virtual button press
}

type InputView interface {
	FocusView

	Input() string
	SetInput(input string)
}

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
func (b baseView) ID() uint64 {
	return b.id
}
func (b *baseView) SetParent(parent View) {
	b.parent = parent
}

type baseFocusView struct {
	focused bool
}

func (b *baseFocusView) SetFocused(focused bool) {
	b.focused = focused
}
func (b *baseFocusView) Focused() bool {
	return b.focused
}
