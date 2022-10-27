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
							ui.Label(ui.DefaultLabelOptions(), "New Game"),
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
	select {
	case id := <-s.buttonPressed:
		switch id {
		case 1: // New Game button
			log.Println("mainMenu - \"New Game\" button pressed")
			manager.Switch(NewNewWorldScene())
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

	// world seed will be received through this channel
	worldSeed chan string
}

func NewNewWorldScene() *newWorldScene {
	worldSeed := make(chan string, 1)

	return &newWorldScene{
		worldSeed: worldSeed,
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow"),
			ui.Center(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 1.0},
				ui.Label(ui.DefaultLabelOptions(), "Enter world seed:"),
				ui.Input(
					func(input string) {
						worldSeed <- input
					},
					ebiten.KeyEnter,
				),
			)))),
	}
}

func (s *newWorldScene) Update(manager *scene.SceneManager) error {
	err := s.view.Update()
	if err != nil {
		return err
	}

	select {
	case seed_string := <-s.worldSeed:
		// convert string to bytes -> compute hash -> convert hash to int64
		seed_bytes := []byte(seed_string)
		seed_hash_bytes := fnv.New64a().Sum(seed_bytes)
		var seed int64
		binary.Read(bytes.NewReader(seed_hash_bytes), binary.BigEndian, &seed)

		log.Printf("newWorldScene - Generating a world with seed %v", seed)

		// Inserting (not pushing) the game scene into the queue
		// So, the queue looks like this
		// [GameScene, MainMenu]
		// And, when the game scene exits, we get back to the main menu.
		manager.Insert(game.New(seed))
		manager.End()
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
