package blocks

import (
	"log"
	"math"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type TexturedBlockState struct {
	Name     string
	Rotation float64
}

// Another base structure, to simplify things
type texturedBlock struct {
	tex      types.Texture
	rotation float64 // in degrees
}

func (b *texturedBlock) Render(_ types.World, screen *ebiten.Image, pos types.Coords2f) {
	opts := &ebiten.DrawImageOptions{}

	if b.rotation != 0 {
		w, h := b.tex.Texture().Size()
		// Move image half a texture size, so that rotation origin will be in the center
		opts.GeoM.Translate(float64(-w/2), float64(-h/2))
		opts.GeoM.Rotate(b.rotation * (math.Pi / 180))
		pos.X += float64(w / 2)
		pos.Y += float64(h / 2)
	}

	opts.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(b.tex.Texture(), opts)
}

func (b *texturedBlock) TextureName() string {
	return b.tex.Name()
}

func (b *texturedBlock) State() interface{} {
	return TexturedBlockState{
		Name:     b.tex.Name(),
		Rotation: b.rotation,
	}
}

func (b *texturedBlock) LoadState(s interface{}) {
	state, ok := s.(TexturedBlockState)
	if !ok {
		log.Panicf("%T - invalid state type; expected %T, got %T", b, TexturedBlockState{}, state)
	}

	b.tex = asset_loader.Texture(state.Name)
	b.rotation = state.Rotation
}
