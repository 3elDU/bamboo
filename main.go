package main

import (
	"math/rand"
	"time"

	"github.com/3elDU/bamboo/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("bamboo devtest")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := game.Create("./assets/")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
