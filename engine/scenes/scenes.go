/*
Various simple scenes, used in the game
*/

package scenes

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"log"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/scene"
	"github.com/3elDU/bamboo/game"
	"github.com/3elDU/bamboo/ui"
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
			ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow"), ui.Padding(1,
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

func (s *mainMenu) Update(manager *scene.SceneManager) error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case id := <-s.buttonPressed:
		switch id {
		case 1: // Singleplayer button
			log.Println("mainMenu - \"Singleplayer\" button pressed")
			manager.Switch(NewWorldListScene())
		case 2: // About
			log.Println("mainMenu - \"About\" button pressed")
			manager.Switch(NewAboutScene())
		case 3: // Exit
			log.Println("mainMenu - \"Exit\" button pressed")
			manager.Exit()
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
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow"),
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

func (s *aboutScene) Update(manager *scene.SceneManager) error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case <-s.goBackEvent:
		manager.End()
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

		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow"), ui.Center(
			ui.Form(
				"Create a new world",
				formData,
				ui.FormPrompt{Title: "World name"},
				ui.FormPrompt{Title: "World seed"},
			),
		))),
	}
}

func (s *newWorldScene) Update(manager *scene.SceneManager) error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case formData := <-s.formData:
		_, seed_string := formData[0], formData[1]

		// convert string to bytes -> compute hash -> convert hash to int64
		seed_bytes := []byte(seed_string)
		seed_hash_bytes := fnv.New64a().Sum(seed_bytes)
		var seed int64
		binary.Read(bytes.NewReader(seed_hash_bytes), binary.BigEndian, &seed)

		manager.QSwitch(game.NewGameScene(seed))
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
	view ui.View

	// when the world will be selected by the user,
	// world name will be transmitted through this channel
	// selectedWorld chan uuid.UUID

	// when the "New world" button will be pressed
	// the event will be transmitted through this channel
	newWorld chan bool
}

func NewWorldListScene() *worldListScene {
	log.Println("NewWorldListScene() - parsing worlds...")

	/*
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
	*/

	view := ui.Stack(ui.StackOptions{
		Direction: ui.VerticalStack,
		Spacing:   1,
	})

	/*
		selectedWorld := make(chan uuid.UUID, 1)
		for _, world := range worldList {
			// extract world UUID here, to use it later in button handler
			worldUUID := world.UUID

			// assemble a view for each world
			view.AddChild(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 0.5},
				ui.Stack(ui.StackOptions{Direction: ui.HorizontalStack, Spacing: 1},
					ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Name: %v", world.Name)),
					ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Seed: %v", world.Seed)),
					ui.Label(ui.DefaultLabelOptions(), fmt.Sprintf("Size: %v", world.Size)),
				),
				ui.Button(func() { selectedWorld <- worldUUID }, ui.Label(ui.DefaultLabelOptions(), "Open world")),
			))
		}
	*/

	newWorld := make(chan bool, 1)
	view.AddChild(ui.Center(
		ui.Button(func() { newWorld <- true }, ui.Label(ui.DefaultLabelOptions(), "New world")),
	))

	return &worldListScene{
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow"), view)),

		// selectedWorld: selectedWorld,
		newWorld: newWorld,
	}
}

func (s *worldListScene) Destroy() {
	log.Println("worldListScene.Destroy() called")
}

func (s *worldListScene) Update(manager *scene.SceneManager) error {
	if err := s.view.Update(); err != nil {
		return err
	}

	select {
	/*
		case id := <-s.selectedWorld:
			log.Printf("worldListScene - Selected world '%v'", id)
			manager.Switch(NewNotImplementedYetScene("World loading"))
	*/
	case <-s.newWorld:
		log.Println("worldListScene - New world")
		manager.QSwitch(NewNewWorldScene())
	default:
	}
	return nil
}

func (s *worldListScene) Draw(screen *ebiten.Image) {
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
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow"),
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

func (s *notImplementedYetScene) Update(manager *scene.SceneManager) error {
	if s.back.IsPressed() {
		manager.End()
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
