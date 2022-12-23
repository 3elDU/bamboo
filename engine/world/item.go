package world

import (
	"github.com/3elDU/bamboo/engine/texture"
	"github.com/hajimehoshi/ebiten/v2"
)

type Item interface {
	Texture() *ebiten.Image
	Type() BlockType
}

type item struct {
	Tex      texture.Texture
	ID       BlockType
	Quantity uint
}

func NewCustomItem(tex texture.Texture, id BlockType, quantity uint) item {
	return item{
		Tex:      tex,
		ID:       id,
		Quantity: quantity,
	}
}

func (i item) Texture() *ebiten.Image {
	return i.Tex.Texture()
}

func (i item) Type() BlockType {
	return i.ID
}
