package texture

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Simple containers that holds both texture name and the texture itself
type Texture struct {
	Name    string
	Texture *ebiten.Image
}
type ConnectedTexture struct {
	Base           string
	SidesConnected [4]bool // in order: left, right, top, bottom
	Texture        *ebiten.Image
}

func (c ConnectedTexture) FullName() string {
	assembledName := c.Base + "-"

	for _, side := range c.SidesConnected {
		if side {
			assembledName += "t"
		} else {
			assembledName += "f"
		}
	}

	return assembledName
}

func (c *ConnectedTexture) SetSides(left, right, top, bottom bool) {
	c.SidesConnected = [4]bool{left, right, top, bottom}
}

func (c *ConnectedTexture) SetSidesArray(sides [4]bool) {
	c.SidesConnected = sides
}
