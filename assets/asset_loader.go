/*
	Loads all assets from the folder
*/

package assets

import (
	_ "image/png"
	"log"
	"path/filepath"
	"strings"

	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type AssetList struct {
	Textures          map[string]*ebiten.Image
	ConnectedTextures map[connectedTexture]*ebiten.Image

	Font *ebiten.Image
}

var (
	GlobalAssets *AssetList = nil
)

// removes file extension, and other parts from the filename
// Example: assets/pictures/picture.png -> picture
func cleanPath(path string) string {
	return strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1)
}

func DefaultFont() *ebiten.Image {
	return GlobalAssets.Font
}

// Texture panicks when a specified texture doesn't exist
func Texture(name string) types.Texture {
	_, exists := GlobalAssets.Textures[name]
	if !exists {
		log.Panicf("texture %v doesn't exist", name)
	}
	return &texture{
		name: name,
	}
}

// ConnectedTexture panicks when a specified texture doesn't exist
func ConnectedTexture(baseName string, left, right, top, bottom bool) types.ConnectedTexture {
	tex := connectedTexture{
		baseName:       baseName,
		connectedSides: [4]bool{left, right, top, bottom},
	}
	_, exists := GlobalAssets.ConnectedTextures[tex]
	if !exists {
		log.Panicf("connected texture %v doesn't exist", tex)
	}
	return &tex
}

// Same as ConnectedTexture, but accepts an array of four booleans
func ConnectedTextureFromArray(baseName string, sidesConnected [4]bool) types.ConnectedTexture {
	return ConnectedTexture(baseName, sidesConnected[0], sidesConnected[1], sidesConnected[2], sidesConnected[3])
}
