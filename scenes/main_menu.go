package scenes

import (
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/scene_manager"
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

func (s *mainMenu) Update() {
	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a view: %v", err)
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
