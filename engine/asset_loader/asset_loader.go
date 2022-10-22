/*
	Loads all assets from the folder
*/

package asset_loader

import (
	"bytes"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/3elDU/bamboo/config"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type AssetList struct {
	Fonts    map[string]font.Face
	Textures map[string]*ebiten.Image

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

// LoadAssets loads assets from directory dir to global variable GlobalAssets
func LoadAssets(dir string) error {
	assetList := &AssetList{
		Fonts:    make(map[string]font.Face),
		Textures: make(map[string]*ebiten.Image),
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			switch filepath.Ext(path) {
			case ".png":
				img, _, err := image.Decode(bytes.NewReader(data))
				if err != nil {
					return err
				}
				assetList.Textures[cleanPath(path)] = ebiten.NewImageFromImage(img)

			case ".ttf":
				f, err := truetype.Parse(data)
				if err != nil {
					return err
				}

				fFace := truetype.NewFace(f, &truetype.Options{
					Size:    config.FontSize,
					Hinting: font.HintingFull,
				})

				assetList.Fonts[strings.Replace(cleanPath(path), "_default", "", 1)] = fFace

				// if filename (without extension) ends in _default, then set this font as default
				if strings.HasSuffix(cleanPath(path), "_default") {
					assetList.defaultFont = fFace
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	GlobalAssets = assetList
	return nil
}

func DefaultFont() font.Face {
	return GlobalAssets.defaultFont
}

// Texture panicks when a specified texture doesn't exist
func Texture(name string) *ebiten.Image {
	tex, exists := GlobalAssets.Textures[name]
	if !exists {
		log.Panicf("texture %v doesn't exist", name)
	}
	return tex
}

// Font panicks when a specified font doesn't exist
func Font(name string) font.Face {
	face, exists := GlobalAssets.Fonts[name]
	if !exists {
		log.Panicf("font %v doesn't exist", name)
	}
	return face
}
