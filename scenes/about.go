package scenes

import (
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/MakeNowJust/heredoc"
	"github.com/hajimehoshi/ebiten/v2"
)

type AboutScene struct {
	view        ui.View
	goBackEvent chan int
}

func NewAboutScene() *AboutScene {
	goBackEvent := make(chan int, 1)

	return &AboutScene{
		goBackEvent: goBackEvent,
		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(),
			ui.Center(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 1},
				ui.Label(ui.DefaultLabelOptions(), heredoc.Doc(`
					Very important text...
					Blah blah blah...
					// TODO: Actually write something here
				`)),
				ui.Button(
					goBackEvent, 1,
					ui.Label(ui.DefaultLabelOptions(), "Back"),
				),
			)),
		)),
	}
}

func (s *AboutScene) Update() {
	err := s.view.Update()
	if err != nil {
		log.Panicf("failed to update a view: %v", err)
	}

	select {
	case <-s.goBackEvent:
		scene_manager.Pop()
	default:
	}
}

func (s *AboutScene) Destroy() {
	log.Println("AboutScene.Destroy() called")
}

func (s *AboutScene) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}
