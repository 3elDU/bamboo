package scenes

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/game"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/3elDU/bamboo/world"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type worldListScene struct {
	worldList []world.WorldSave
	view      ui.View

	// when the world will be selected by the user,
	// world name will be transmitted through this channel
	selectedWorld chan uuid.UUID
	deleteWorld   chan uuid.UUID

	// when the "New world" button will be pressed
	// the event will be transmitted through this channel
	newWorld chan bool
}

// Scans the save folder for worlds
func (s *worldListScene) Scan() {
	log.Printf("worldListScene.Scan()")

	worldList := make([]world.WorldSave, 0)
	filepath.WalkDir(config.WorldSaveDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// only check for directories
		if !d.IsDir() {
			return nil
		}

		// read world metadata
		worldInfo, err := os.Open(filepath.Join(path, "world.gob"))
		if err != nil {
			// skip the directory, if it doesn't have the "world.gob" file inside of it,
			// or if it is inaccesible for some other reason
			// but don't throw an error!
			return nil
		}
		defer worldInfo.Close()

		decoder := gob.NewDecoder(worldInfo)
		worldMetadata := new(world.WorldSave)
		if err = decoder.Decode(worldMetadata); err != nil {
			return err
		}

		log.Println(*worldMetadata)
		worldList = append(worldList, *worldMetadata)

		return nil
	})
	s.worldList = worldList
}

func (s *worldListScene) UpdateUI() {
	view := ui.Stack(ui.StackOptions{
		Direction: ui.VerticalStack,
		Spacing:   1,
	})
	for _, world := range s.worldList {
		// extract world UUID here, to use it later in button handler
		worldUUID := world.UUID

		// assemble a view for each world
		view.AddChild(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 0.5},
			ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 1},
				ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Name: %v", world.Name)),
				ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Seed: %v", world.Seed)),
				ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Size: %v", world.Size)),
			),
			ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 1},
				ui.Button(func() { s.selectedWorld <- worldUUID }, ui.Label(ui.DefaultLabelOptions(), "Play")),
				ui.Button(func() { s.deleteWorld <- worldUUID }, ui.Label(ui.DefaultLabelOptions(), "Delete")),
			),
		))
	}
	view.AddChild(ui.Center(
		ui.Button(func() { s.newWorld <- true }, ui.Label(ui.DefaultLabelOptions(), "New world")),
	))

	s.view = ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(), view))
}

func NewWorldListScene() *worldListScene {
	scene := &worldListScene{
		selectedWorld: make(chan uuid.UUID, 1),
		deleteWorld:   make(chan uuid.UUID, 1),
		newWorld:      make(chan bool, 1),
	}
	scene.Scan()
	return scene
}

func (s *worldListScene) Destroy() {
	log.Println("worldListScene.Destroy() called")
}

func (s *worldListScene) Update() {
	// Rescan the saves folder each 60 ticks ( 1 second )
	if scene_manager.Ticks()%60 == 0 {
		s.Scan()
	}

	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a view: %v", err)
	}

	select {
	case id := <-s.selectedWorld:
		log.Printf("worldListScene - Selected world '%v'", id)
		loadedWorld := world.LoadWorld(id)
		loadedPlayer := player.LoadPlayer(id)
		scene_manager.QSwitch(game.NewGameScene(loadedWorld, *loadedPlayer))
	case <-s.newWorld:
		log.Println("worldListScene - New world")
		scene_manager.QSwitch(NewNewWorldScene())
	case id := <-s.deleteWorld:
		world.DeleteWorld(id)
		s.Scan()
	default:
	}
}

func (s *worldListScene) Draw(screen *ebiten.Image) {
	if scene_manager.Ticks()%60 == 0 || s.view == nil {
		s.UpdateUI()
	}
	if err := s.view.Draw(screen, 0, 0); err != nil {
		log.Panicf("worldListScene.view.Draw() - %v", err)
	}
}
