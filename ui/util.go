package ui

import "github.com/hajimehoshi/ebiten/v2"

// A simple wrapper that initializes the screen component and immediately draws it
func ImmediateDraw(screen *ebiten.Image, child Component) {
	Screen(child).Draw(screen, 0, 0)
}
