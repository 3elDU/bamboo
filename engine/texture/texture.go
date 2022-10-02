/*
	Some helper methods on textures for easier querying, modification, etc.
*/

package texture

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Panicks if call to Texture.Query() fails
func Width(t *sdl.Texture) int32 {
	_, _, width, _, err := t.Query()

	if err != nil {
		panic(err)
	}

	return width
}

// Panicks if call to Texture.Query() fails
func Height(t *sdl.Texture) int32 {
	_, _, _, height, err := t.Query()
	if err != nil {
		panic(err)
	}
	return height
}

// Returns both width and height
// Panicks if call to Texture.Query() fails
func Dimensions(t *sdl.Texture) (int32, int32) {
	_, _, width, height, err := t.Query()
	if err != nil {
		panic(err)
	}
	return width, height
}
