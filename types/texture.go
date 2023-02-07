package types

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Texture interface {
	Texture() *ebiten.Image
	Name() string
}

type ConnectedTexture interface {
	Texture() *ebiten.Image
	ConnectedSides() [4]bool
	SetConnectedSides(sides [4]bool)
	Name() string
}
