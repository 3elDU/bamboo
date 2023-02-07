package asset_loader

import (
	"bytes"
	"image"
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

func parseTexture(assetList *AssetList, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	assetList.Textures[cleanPath(path)] = ebiten.NewImageFromImage(img)

	return nil
}

func parseConnectedTexture(assetList *AssetList, path string) error {
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

	// texture map describes, which sub-texture is on which coordinate
	// first and second indices represent coordinates ( multiples of 16 ) on an atlas
	// the [4]bool array describes the connected texture itself
	// sides go in order: left, right, top, bottom
	textureMap := [4][4][4]bool{
		{{false, false, false, false}, {false, true, false, false}, {true, true, false, false}, {true, false, false, false}},
		{{false, false, false, true}, {false, true, false, true}, {true, true, false, true}, {true, false, false, true}},
		{{false, false, true, true}, {false, true, true, true}, {true, true, true, true}, {true, false, true, true}},
		{{false, false, true, false}, {false, true, true, false}, {true, true, true, false}, {true, false, true, false}},
	}

	for y, row := range textureMap {
		for x, col := range row {
			assetList.ConnectedTextures[connectedTexture{
				baseName:       cleanPath(path),
				connectedSides: col,
			}] = ebiten.NewImageFromImage(
				tex.SubImage(image.Rect(x*16, y*16, x*16+16, y*16+16)),
			)
		}
	}

	return nil
}

func parseFont(assetList *AssetList, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

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

	return nil
}

// LoadAssets loads assets from directory dir to global variable GlobalAssets
func LoadAssets(dir string) {
	assetList := &AssetList{
		Fonts:             make(map[string]font.Face),
		Textures:          make(map[string]*ebiten.Image),
		ConnectedTextures: make(map[connectedTexture]*ebiten.Image),
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If there is a directory, treat it as a texture atlas
		if d.IsDir() {
			return parseConnectedTexture(assetList, path)
		}

		switch filepath.Ext(path) {
		case ".png":
			return parseTexture(assetList, path)
		case ".ttf":
			return parseFont(assetList, path)
		}

		return nil
	})
	if err != nil {
		log.Panicln(err)
	}

	GlobalAssets = assetList
}
