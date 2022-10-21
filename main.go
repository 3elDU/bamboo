package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	if err := asset_loader.LoadAssets(config.AssetDirectory); err != nil {
		panic(fmt.Sprintf("LoadAssets() failed with: %v", err))
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("bamboo devtest")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.Create()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
