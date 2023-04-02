package asset_loader

import (
	"github.com/3elDU/bamboo/config"
	"github.com/hajimehoshi/ebiten/v2"
)

type texture struct {
	name string
}

func (t *texture) Texture() *ebiten.Image {
	return GlobalAssets.Textures[t.name]
}

func (t *texture) Name() string {
	return t.name
}

func (t *texture) ScaledSize() (float64, float64) {
	w, h := t.Texture().Size()
	return float64(w) * config.UIScaling, float64(h) * config.UIScaling
}

type connectedTexture struct {
	baseName       string
	connectedSides [4]bool
}

func (t *connectedTexture) Texture() *ebiten.Image {
	return GlobalAssets.ConnectedTextures[*t]
}

func (t *connectedTexture) ConnectedSides() [4]bool {
	return t.connectedSides
}

func (t *connectedTexture) SetConnectedSides(sides [4]bool) {
	t.connectedSides = sides
}

func (t *connectedTexture) Name() string {
	return t.baseName
}
