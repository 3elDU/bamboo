package scenes

import (
	"log"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ConfirmationScene struct {
	prompt string
	view   ui.Component

	confirmAction func()
	noButton      chan bool
	yesButton     chan bool
}

// Asks user to confirm action with "Yes" or "No"
func NewConfirmationScene(prompt string, confirmAction func()) scene_manager.Scene {
	noButton := make(chan bool, 1)
	yesButton := make(chan bool, 1)

	view := ui.Screen(ui.BackgroundColor(colors.C("white"), ui.Center(
		ui.VStack().WithSpacing(2.0).WithChildren(
			ui.HCenter(ui.Label(prompt)),
			ui.HCenter(ui.HStack().WithSpacing(1.0).WithChildren(
				ui.Button(noButton, true, ui.Label("No")),
				ui.Button(yesButton, true, ui.Label("Yes")),
			)),
		),
	)))

	return &ConfirmationScene{
		prompt: prompt,
		view:   view,

		confirmAction: confirmAction,
		noButton:      noButton,
		yesButton:     yesButton,
	}
}

func (s *ConfirmationScene) Update() {
	if err := s.view.Update(); err != nil {
		log.Panic(err)
	}

	select {
	case <-s.noButton:
		scene_manager.Pop()
	case <-s.yesButton:
		s.confirmAction()
		scene_manager.Pop()
	default:
	}
}

func (s *ConfirmationScene) Draw(screen *ebiten.Image) {
	if err := s.view.Draw(screen, 0, 0); err != nil {
		log.Panic(err)
	}
}

func (s *ConfirmationScene) Destroy() {
	log.Println("ConfirmationScene.Destroy() called")
}
