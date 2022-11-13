package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/scene_manager"
	"github.com/3elDU/bamboo/engine/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	// init logging
	logFilename := "Log-" + time.Now().Format("02-Jan-2006_15-04-05-MST") + ".txt"
	file, err := os.Create(filepath.Join("logs", logFilename))
	if err != nil {
		panic(fmt.Sprintf("failed to create log file: %v", err))
	}
	w := io.MultiWriter(file, os.Stdout)
	log.SetOutput(w)

	// init RNG
	rand.Seed(int64(time.Now().Nanosecond()))

	// load assets
	if err := asset_loader.LoadAssets(config.AssetDirectory); err != nil {
		log.Panicf("LoadAssets() failed with %v", err)
	}

	// set window options
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("bamboo devtest")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func main() {
	// init scene manager, and scenes
	scene_manager.Push(scenes.NewMainMenuScene())

	// run main loop!
	scene_manager.Run()
}
