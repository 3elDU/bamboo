package types

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Texture interface {
	Texture() *ebiten.Image
	Name() string
	// Returns size of the texture, multiplied by ui scaling
	// Useful for UI elements
	ScaledSize() (float64, float64)
}

type ConnectedTexture interface {
	Texture() *ebiten.Image
	ConnectedSides() [4]bool
	SetConnectedSides(sides [4]bool)
	Name() string
}
