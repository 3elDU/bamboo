package scenes

import (
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

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

func (s *aboutScene) Update() {
	err := s.view.Update()
	if err != nil {
		log.Panicf("failed to update a view: %v", err)
	}

	select {
	case <-s.goBackEvent:
		scene_manager.End()
	default:
	}
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
