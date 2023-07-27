package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/profile"
	"golang.org/x/exp/slices"

	// imports for side effects
	_ "github.com/3elDU/bamboo/blocks_impl"
	_ "github.com/3elDU/bamboo/items_impl"
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

	// load assets
	asset_loader.LoadAssets(config.AssetDirectory)

	// set window options
	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("bamboo devtest")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func main() {
	if slices.Contains(os.Environ(), "CPUPROFILE=1") {
		log.Println("Starting with CPU profiling enabled")
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}
	if slices.Contains(os.Environ(), "MEMPROFILE=1") {
		log.Println("Starting with memory profiling enabled")
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	}

	// init scene manager, and scenes
	scene_manager.Push(scenes.NewMainMenuScene())

	// run main loop!
	scene_manager.Run()
}
