package scenes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world_type"

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
	view      ui.Component

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
func (scene *WorldListScene) Scan() {
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
		worldInfo, err := os.ReadFile(filepath.Join(path, "world.gob"))
		if err != nil {
			// skip the directory, if it doesn't have the "world.gob" file inside of it,
			// or if it is inaccesible for some other reason
			// but don't throw an error!
			return nil
		}

		decoder := gob.NewDecoder(bytes.NewReader(worldInfo))
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
	scene.worldList = worldList
}

func (scene *WorldListScene) UpdateUI() {
	rootView := ui.VStack().WithSpacing(3)
	worldList := ui.VStack().WithSpacing(1.0)

	for _, currentWorld := range scene.worldList {
		// assemble a view for each world
		worldList.AddChild(ui.VStack().WithSpacing(0.5).AlignChildren(ui.AlignCenter).WithChildren(
			ui.HStack().WithSpacing(0.5).WithChildren(
				ui.Label(fmt.Sprintf("Name: %v", currentWorld.Name)),
				ui.Label(fmt.Sprintf("Seed: %v", currentWorld.Seed)),
			),
			ui.HStack().WithSpacing(1).WithChildren(
				ui.Button(scene.selectedWorld, currentWorld, ui.Label(fmt.Sprintf("Play %v", currentWorld.Name))),
				ui.Button(scene.deleteWorld, currentWorld, ui.Label("Delete")),
			),
		))
	}
	rootView.AddChild(worldList)

	rootView.AddChild(ui.HStack().WithSpacing(1).WithChildren(
		ui.Button(scene.newWorld, true, ui.Label("New world")),
		ui.Button(scene.goBack, true, ui.Label("Go back")),
	))

	scene.view = ui.Screen(
		ui.TileBackgroundImage(asset_loader.Texture("snow"),
			ui.Center(rootView),
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

func (scene *WorldListScene) Destroy() {
	log.Println("worldListScene.Destroy() called")
}

func (scene *WorldListScene) Update() {
	// Rescan the saves folder each 60 ticks ( 1 second )
	if scene_manager.Ticks()%60 == 0 {
		scene.Scan()
		scene.UpdateUI()
	}

	if err := scene.view.Update(); err != nil {
		log.Panicf("failed to update a view: %v", err)
	}

	select {
	case save := <-scene.selectedWorld:
		log.Printf("worldListScene - Selected world '%v'", save)
		scene_manager.ReplaceAndSwitch(game.LoadGameScene(save))
	case <-scene.newWorld:
		log.Println("worldListScene - New world")
		scene_manager.PushAndSwitch(NewNewWorldScene())
	case <-scene.goBack:
		scene_manager.Pop()
	case metadata := <-scene.deleteWorld:
		scene_manager.PushAndSwitch(NewConfirmationScene(
			fmt.Sprintf("Are you sure to do delete the world \"%v\"?", metadata.Name),
			func() {
				world.DeleteWorld(metadata)
				scene.Scan()
				scene.UpdateUI()
			},
		))
	default:
	}
}

func (scene *WorldListScene) Draw(screen *ebiten.Image) {
	if scene_manager.Ticks()%60 == 0 || scene.view == nil {
		scene.UpdateUI()
	}
	if err := scene.view.Draw(screen, 0, 0); err != nil {
		log.Panicf("worldListScene.view.Draw() - %v", err)
	}
}
