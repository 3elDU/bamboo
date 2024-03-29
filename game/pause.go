// Pause menu

package game

import (
	"image/color"
	"log"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type pauseMenu struct {
	view ui.Component

	// just a black texture
	// it is used to "dim" the background
	tex  *ebiten.Image
	opts *ebiten.DrawImageOptions

	// button press event will be received through these channels
	continueBtn, exitBtn chan bool
}

func newPauseMenu() *pauseMenu {
	tex := ebiten.NewImage(1, 1)
	tex.Fill(color.RGBA{A: 128})

	var (
		continueBtn = make(chan bool, 1)
		exitBtn     = make(chan bool, 1)
	)

	return &pauseMenu{
		tex:  tex,
		opts: &ebiten.DrawImageOptions{},

		continueBtn: continueBtn,
		exitBtn:     exitBtn,

		view: ui.Screen(
			ui.Padding(1,
				ui.VStack().
					WithProportions(0.3).
					WithChildren(
						ui.Center(
							ui.CustomLabel("Paused", colors.C("white"), 3.0),
						),

						ui.Center(ui.VStack().WithSpacing(1.0).WithChildren(
							ui.Button(continueBtn, true, ui.Label("Continue")),
							ui.Button(exitBtn, true, ui.Label("Exit to main menu")),
						)),
					),
			),
		),
	}
}

func (p *pauseMenu) Draw(screen *ebiten.Image) {
	// Dim the background with translucent black texture
	p.opts.GeoM.Reset()
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	// Fixes issue #3
	p.opts.GeoM.Scale(float64(w)+2, float64(h)+2)
	screen.DrawImage(p.tex, p.opts)

	if err := p.view.Draw(screen, 0, 0); err != nil {
		log.Panic(err)
	}
}

func (p *pauseMenu) Update() {
	if err := p.view.Update(); err != nil {
		log.Panicf("pauseMenu.ButtonPressed() - %v", err)
	}

	select {
	case <-p.continueBtn:
		log.Println("pauseMenu - \"Continue\" button pressed")
		scene_manager.HideOverlay()
	case <-p.exitBtn:
		log.Println("pauseMenu - \"Exit to main menu\" button pressed")
		scene_manager.HideOverlay()
		scene_manager.Pop()
	default:
	}
}

func (p *pauseMenu) Destroy() {

}
