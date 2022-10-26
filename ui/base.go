/*
	Base structures and interfaces
*/

package ui

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type View interface {
	SetParent(parent View)
	Children() []View

	// maximum space that the component could theoretically occupy in their container
	MaxSize() (float64, float64)
	// how much space the component takes in practice
	ComputedSize() (float64, float64)
	// practical maximum space, that provided child could take in the container
	CapacityForChild(child View) (float64, float64)

	// There is no Update() function
	// Update of the component logic shall happen in the Draw() function
	Draw(screen *ebiten.Image, x, y float64) error

	// returns unique identifier of the component, so it can be compared to others
	ID() uint64
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
