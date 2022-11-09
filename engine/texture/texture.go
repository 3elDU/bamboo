package texture

import "github.com/hajimehoshi/ebiten/v2"

// Simple containers that holds both texture name and the texture itself
type Texture struct {
	Name    string
	Texture *ebiten.Image
}
