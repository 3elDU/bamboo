/*
Various simple scenes, used in the game
*/

package scenes

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/scene_manager"
	"github.com/3elDU/bamboo/engine/world"
	"github.com/3elDU/bamboo/game"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/ui"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type mainMenu struct {
	view ui.View

	// through this channel we will receive button ID, that was pressed
	buttonPressed chan int
}

func NewMainMenuScene() *mainMenu {
	buttonPressed := make(chan int, 1)

	return &mainMenu{
		buttonPressed: buttonPressed,
		view: ui.Screen(
			ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(), ui.Padding(1,
				ui.Stack(ui.StackOptions{Spacing: 0, Proportions: []float64{0.2}},
					ui.Center(ui.Label(
						ui.LabelOptions{
							Color:   colors.Black,
							Scaling: 2.5,
						},
						"bamboo devtest",
					)),
					ui.Center(ui.Stack(ui.StackOptions{Spacing: 0.5},
						ui.Button(
							func() { buttonPressed <- 1 },
							ui.Label(ui.DefaultLabelOptions(), "Singleplayer"),
						),
						ui.Button(
							func() { buttonPressed <- 2 },
							ui.Label(ui.DefaultLabelOptions(), "About"),
						),
						ui.Button(
							func() { buttonPressed <- 3 },
							ui.Label(ui.DefaultLabelOptions(), "Exit"),
						),
					)),
				),
			)),
		),
	}
}

func (s *mainMenu) Update() error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case id := <-s.buttonPressed:
		switch id {
		case 1: // Singleplayer button
			log.Println("mainMenu - \"Singleplayer\" button pressed")
			scene_manager.Switch(NewWorldListScene())
		case 2: // About
			log.Println("mainMenu - \"About\" button pressed")
			scene_manager.Switch(NewAboutScene())
		case 3: // Exit
			log.Println("mainMenu - \"Exit\" button pressed")
			scene_manager.Exit()
		}
	default:
	}

	return nil
}

func (s *mainMenu) Destroy() {
	log.Println("mainMenu.Destroy() called")
}

func (s *mainMenu) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}

type aboutScene struct {
	view        ui.View
	goBackEvent chan int
}

func NewAboutScene() *aboutScene {
	goBackEvent := make(chan int, 1)

	return &aboutScene{
		goBackEvent: goBackEvent,
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(),
			ui.Center(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 1},
				ui.Label(ui.DefaultLabelOptions(), "Very important text..."),
				ui.Label(ui.DefaultLabelOptions(), "Blah blah blah..."),
				ui.Label(ui.DefaultLabelOptions(), "// TODO: Actually write something here"),
				ui.Button(
					func() { goBackEvent <- 1 },
					ui.Label(ui.DefaultLabelOptions(), "Back"),
				),
			)),
		)),
	}
}

func (s *aboutScene) Update() error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case <-s.goBackEvent:
		scene_manager.End()
	default:
	}
	return nil
}

func (s *aboutScene) Destroy() {
	log.Println("aboutScene.Destroy() called")
}

func (s *aboutScene) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}

type newWorldScene struct {
	view ui.View

	// form results will be received through this channel
	// first string is world name, second is world seed
	formData chan []string
}

func NewNewWorldScene() *newWorldScene {
	formData := make(chan []string, 1)

	return &newWorldScene{
		formData: formData,

		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(), ui.Center(
			ui.Form(
				"Create a new world",
				formData,
				ui.FormPrompt{Title: "World name"},
				ui.FormPrompt{Title: "World seed"},
			),
		))),
	}
}

func (s *newWorldScene) Update() error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case formData := <-s.formData:
		world_name, seed_string := formData[0], formData[1]

		// convert string to bytes -> compute hash -> convert hash to int64
		seed_bytes := []byte(seed_string)
		seed_hash_bytes := fnv.New64a().Sum(seed_bytes)
		var seed int64
		binary.Read(bytes.NewReader(seed_hash_bytes), binary.BigEndian, &seed)

		w := world.NewWorld(world_name, uuid.New(), seed)
		scene_manager.QSwitch(game.NewGameScene(w, player.Player{X: config.PlayerStartX, Y: config.PlayerStartY}))
	default:
	}
	return nil
}

func (*newWorldScene) Destroy() {
	log.Println("newWorldScene.Destroy() called")
}

func (s *newWorldScene) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}

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

func (s *worldListScene) Update() error {
	// Rescan the saves folder each 60 ticks ( 1 second )
	if scene_manager.Ticks()%60 == 0 {
		s.Scan()
	}

	if err := s.view.Update(); err != nil {
		return err
	}

	select {
	case id := <-s.selectedWorld:
		log.Printf("worldListScene - Selected world '%v'", id)
		w, p := world.LoadWorld(id)
		scene_manager.QSwitch(game.NewGameScene(w, p))
	case <-s.newWorld:
		log.Println("worldListScene - New world")
		scene_manager.QSwitch(NewNewWorldScene())
	case id := <-s.deleteWorld:
		world.DeleteWorld(id)
		s.Scan()
	default:
	}
	return nil
}

func (s *worldListScene) Draw(screen *ebiten.Image) {
	if scene_manager.Ticks()%60 == 0 || s.view == nil {
		s.UpdateUI()
	}
	if err := s.view.Draw(screen, 0, 0); err != nil {
		log.Panicf("worldListScene.view.Draw() - %v", err)
	}
}

type notImplementedYetScene struct {
	view ui.View
	back ui.ButtonView
}

func NewNotImplementedYetScene(thing string) *notImplementedYetScene {
	backButton := ui.Button(func() {}, ui.Label(ui.DefaultLabelOptions(), "Back"))

	return &notImplementedYetScene{
		back: backButton,
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(),
			ui.Stack(ui.StackOptions{Direction: ui.VerticalStack},
				ui.Center(ui.Label(ui.DefaultLabelOptions(), thing+" isn't implemented yet!")),
				ui.Center(backButton),
			),
		)),
	}
}

func (s *notImplementedYetScene) Destroy() {
	log.Println("notImplementedYetScene.Destroy() called")
}

func (s *notImplementedYetScene) Update() error {
	if s.back.IsPressed() {
		scene_manager.End()
	}

	if err := s.view.Update(); err != nil {
		return err
	}

	return nil
}

func (s *notImplementedYetScene) Draw(screen *ebiten.Image) {
	if err := s.view.Draw(screen, 0, 0); err != nil {
		log.Panicf("notImplementedYetScene.view.Draw() - %v", err)
	}
}
