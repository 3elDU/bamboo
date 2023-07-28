package scenes

import (
	"log"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MainMenu struct {
	view ui.Component

	// through this channel we will receive button ID, that was pressed
	buttonPressed chan int
}

func NewMainMenuScene() *MainMenu {
	buttonPressed := make(chan int, 1)

	return &MainMenu{
		buttonPressed: buttonPressed,
		view: ui.Screen(ui.TileBackgroundImage(assets.Texture("snow"), ui.Padding(0.5,
			ui.Overlay(
				ui.VStack().WithProportions(0.4).WithChildren(
					ui.Center(
						ui.Label("Bamboo").WithTextSize(5),
					),
					ui.Center(ui.VStack().WithSpacing(0.5).WithChildren(
						ui.Button(buttonPressed, 1,
							ui.Label("Singleplayer"),
						),
						ui.Button(buttonPressed, 2,
							ui.Label("About"),
						),
						ui.Button(buttonPressed, 3,
							ui.Label("Exit"),
						),
					)),
				),

				ui.PositionSelf(ui.PositionBottomRight,
					ui.HStack(
						ui.VStack().WithSpacing(0.2).AlignChildren(ui.AlignEnd).WithChildren(
							ui.Label("Version: "),
							ui.Label("Commit hash: "),
							ui.Label("Build date: "),
							ui.Label("Build machine: "),
						),

						ui.VStack().WithSpacing(0.2).AlignChildren(ui.AlignStart).WithChildren(
							ui.Label(config.GitTag),
							ui.Label(config.GitCommit),
							ui.Label(config.BuildDate),
							ui.Label(config.BuildMachine),
						),
					),
				),
			),
		)),
		),
	}
}

func (s *MainMenu) Update() {
	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a view: %v", err)
	}

	// Open debug menu on "D"
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		scene_manager.PushAndSwitch(NewDebugMenu())
	}

	select {
	case id := <-s.buttonPressed:
		switch id {
		case 1: // Singleplayer button
			log.Println("mainMenu - \"Singleplayer\" button pressed")
			scene_manager.PushAndSwitch(NewWorldListScene())
		case 2: // About
			log.Println("mainMenu - \"About\" button pressed")
			scene_manager.PushAndSwitch(NewAboutScene())
		case 3: // Exit
			log.Println("mainMenu - \"Exit\" button pressed")
			scene_manager.Exit()
		}
	default:
	}
}

func (s *MainMenu) Destroy() {
	log.Println("MainMenu.Destroy() called")
}

func (s *MainMenu) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}
