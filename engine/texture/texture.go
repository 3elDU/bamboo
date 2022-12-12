package texture

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Simple containers that holds both texture name and the texture itself
type Texture struct {
	Name    string
	texture *ebiten.Image
}

func NewTexture(name string, texture *ebiten.Image) Texture {
	return Texture{Name: name, texture: texture}
}

func (t Texture) Texture() *ebiten.Image {
	return t.texture
}

type ConnectedTexture struct {
	Base           string
	SidesConnected [4]bool       // in order: left, right, top, bottom
	texture        *ebiten.Image // texture is unexported, so that it won't be exported in world saves
}

func NewConnectedTexture(base string, sidesConnected [4]bool, texture *ebiten.Image) ConnectedTexture {
	return ConnectedTexture{Base: base, SidesConnected: sidesConnected, texture: texture}
}

func (c ConnectedTexture) Texture() *ebiten.Image {
	return c.texture
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
