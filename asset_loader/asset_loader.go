/*
	Loads all assets from the folder
*/

package asset_loader

import (
	_ "image/png"
	"log"
	"path/filepath"
	"strings"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

func init() {
	LoadAssets(config.AssetDirectory)
}

type AssetList struct {
	Fonts             map[string]font.Face
	Textures          map[string]*ebiten.Image
	ConnectedTextures map[connectedTexture]*ebiten.Image

	defaultFont font.Face
}

var (
	GlobalAssets *AssetList = nil
)

// removes file extension, and other parts from the filename
// Example: assets/pictures/picture.png -> picture
func cleanPath(path string) string {
	return strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1)
}

func DefaultFont() font.Face {
	return GlobalAssets.defaultFont
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

// Font panicks when a specified font doesn't exist
func Font(name string) font.Face {
	face, exists := GlobalAssets.Fonts[name]
	if !exists {
		log.Panicf("font %v doesn't exist", name)
	}
	return face
}
