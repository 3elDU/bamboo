package world

import (
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type Item interface {
	Texture() *ebiten.Image
	Type() types.BlockType
}

type item struct {
	Tex      types.Texture
	ID       types.BlockType
	Quantity uint
}

func NewCustomItem(tex types.Texture, id types.BlockType, quantity uint) item {
	return item{
		Tex:      tex,
		ID:       id,
		Quantity: quantity,
	}
}

func (i item) Texture() *ebiten.Image {
	return i.Tex.Texture()
}

func (i item) Type() types.BlockType {
	return i.ID
}
