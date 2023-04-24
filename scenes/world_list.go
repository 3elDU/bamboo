package scenes

import (
	"encoding/gob"
	"fmt"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world_type"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/game"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/3elDU/bamboo/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type WorldListScene struct {
	worldList []types.Save
	view      ui.View

	// when the world will be selected by the user,
	// world name will be transmitted through this channel
	selectedWorld chan types.Save
	deleteWorld   chan types.Save

	// when the "New world" button will be pressed
	// the event will be transmitted through this channel
	newWorld chan bool
	// same idead as for newWorld
	goBack chan bool
}

// Scan scans the save folder for worlds
func (s *WorldListScene) Scan() {
	worldList := make([]types.Save, 0)
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
		worldMetadata := new(types.Save)
		if err = decoder.Decode(worldMetadata); err != nil {
			return err
		}

		if worldMetadata.WorldType != world_type.Overworld {
			return nil
		}

		worldList = append(worldList, *worldMetadata)

		return nil
	})
	s.worldList = worldList
}

func (s *WorldListScene) UpdateUI() {
	view := ui.Stack(ui.StackOptions{
		Direction: ui.VerticalStack,
		Spacing:   3,
	})

	worldList := ui.Stack(ui.StackOptions{
		Direction: ui.VerticalStack,
		Spacing:   1,
	})
	for _, currentWorld := range s.worldList {
		// assemble a view for each world
		worldList.AddChild(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 0.5},
			ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 1},
				ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Name: %v", currentWorld.Name)),
				ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Seed: %v", currentWorld.Seed)),
			),
			ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 1},
				ui.Button(func() { s.selectedWorld <- currentWorld }, ui.Label(ui.DefaultLabelOptions(), "Play")),
				ui.Button(func() { s.deleteWorld <- currentWorld }, ui.Label(ui.DefaultLabelOptions(), "Delete")),
			),
		))
	}
	view.AddChild(worldList)

	view.AddChild(ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 1},
		ui.Button(func() { s.newWorld <- true }, ui.Label(ui.DefaultLabelOptions(), "New world")),
		ui.Button(func() { s.goBack <- true }, ui.Label(ui.DefaultLabelOptions(), "Go back")),
	))

	s.view = ui.Screen(
		ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(),
			ui.Center(view),
		),
	)
}

func NewWorldListScene() *WorldListScene {
	scene := &WorldListScene{
		selectedWorld: make(chan types.Save, 1),
		deleteWorld:   make(chan types.Save, 1),
		newWorld:      make(chan bool, 1),
		goBack:        make(chan bool, 1),
	}
	scene.Scan()
	scene.UpdateUI()
	return scene
}

func (s *WorldListScene) Destroy() {
	log.Println("worldListScene.Destroy() called")
}

func (s *WorldListScene) Update() {
	// Rescan the saves folder each 60 ticks ( 1 second )
	if scene_manager.Ticks()%60 == 0 {
		s.Scan()
	}

	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a view: %v", err)
	}

	select {
	case save := <-s.selectedWorld:
		log.Printf("worldListScene - Selected world '%v'", save)
		scene_manager.QPushAndSwitch(game.LoadGameScene(save))
	case <-s.newWorld:
		log.Println("worldListScene - New world")
		scene_manager.PushAndSwitch(NewNewWorldScene())
	case <-s.goBack:
		scene_manager.Pop()
	case id := <-s.deleteWorld:
		world.DeleteWorld(id)
		s.Scan()
	default:
	}
}

func (s *WorldListScene) Draw(screen *ebiten.Image) {
	if scene_manager.Ticks()%60 == 0 || s.view == nil {
		s.UpdateUI()
	}
	if err := s.view.Draw(screen, 0, 0); err != nil {
		log.Panicf("worldListScene.view.Draw() - %v", err)
	}
}
