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
	"github.com/3elDU/bamboo/engine/texture"
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

		// If there is a directory, treat it as a texture atlas
		if d.IsDir() {
			data, err := os.ReadFile(filepath.Join(path, "atlas.png"))
			if err != nil {
				// Ignore folders without a texture atlas
				return nil
			}

			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				return err
			}

			tex := ebiten.NewImageFromImage(img)

			// extract sub-textures from atlas
			// each character after the texture name resembles a side
			// if a side is connected, it's t, else it's f
			// sides are in order: left, right, top, bottom
			assetList.Textures[d.Name()+"-ffff"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(0, 0, 16, 16)))
			assetList.Textures[d.Name()+"-ftff"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(16, 0, 32, 32)))
			assetList.Textures[d.Name()+"-ttff"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(32, 0, 48, 48)))
			assetList.Textures[d.Name()+"-tfff"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(48, 0, 64, 64)))
			assetList.Textures[d.Name()+"-ffft"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(0, 16, 16, 32)))
			assetList.Textures[d.Name()+"-ftft"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(16, 16, 32, 32)))
			assetList.Textures[d.Name()+"-ttft"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(32, 16, 48, 32)))
			assetList.Textures[d.Name()+"-tfft"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(48, 16, 64, 32)))
			assetList.Textures[d.Name()+"-fftt"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(0, 32, 16, 48)))
			assetList.Textures[d.Name()+"-fttt"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(16, 32, 32, 48)))
			assetList.Textures[d.Name()+"-tttt"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(32, 32, 48, 48)))
			assetList.Textures[d.Name()+"-tftt"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(48, 32, 64, 48)))
			assetList.Textures[d.Name()+"-fftf"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(0, 48, 16, 64)))
			assetList.Textures[d.Name()+"-fttf"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(16, 48, 32, 64)))
			assetList.Textures[d.Name()+"-tttf"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(32, 48, 48, 64)))
			assetList.Textures[d.Name()+"-tftf"] = ebiten.NewImageFromImage(tex.SubImage(image.Rect(48, 48, 64, 64)))
		} else {
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
func Texture(name string) texture.Texture {
	tex, exists := GlobalAssets.Textures[name]
	if !exists {
		log.Panicf("texture %v doesn't exist", name)
	}
	return texture.Texture{
		Name:    name,
		Texture: tex,
	}
}

// ConnectedTexture panicks when a specified texture doesn't exist
func ConnectedTexture(baseName string, left, right, top, bottom bool) texture.ConnectedTexture {
	assembledName := baseName + "-"

	for _, side := range [4]bool{left, right, top, bottom} {
		if side {
			assembledName += "t"
		} else {
			assembledName += "f"
		}
	}

	tex, exists := GlobalAssets.Textures[assembledName]
	if !exists {
		log.Panicf("connected texture %v doesn't exist", assembledName)
	}

	return texture.ConnectedTexture{
		Base:           baseName,
		SidesConnected: [4]bool{left, right, top, bottom},
		Texture:        tex,
	}
}

// Font panicks when a specified font doesn't exist
func Font(name string) font.Face {
	face, exists := GlobalAssets.Fonts[name]
	if !exists {
		log.Panicf("font %v doesn't exist", name)
	}
	return face
}
