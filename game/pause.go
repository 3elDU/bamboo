// Pause menu

package game

import (
	"image/color"
	"log"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type buttonPressedEvent int

const (
	noButtonPressed buttonPressedEvent = iota
	continueButtonPressed
	exitButtonPressed
)

// Pause menu is not a scene
// It is displayed ON TOP of existing game scene
type pauseMenu struct {
	view ui.View

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

		view: ui.Screen(ui.Padding(1, ui.Stack(
			ui.StackOptions{
				Direction:   ui.VerticalStack,
				Proportions: []float64{0.3},
			},

			ui.Center(ui.Label(ui.LabelOptions{Color: colors.White, Scaling: 3.0}, "Paused")),

			ui.Center(ui.Stack(ui.StackOptions{Direction: ui.VerticalStack, Spacing: 1},
				ui.Button(func() { continueBtn <- true }, ui.Label(ui.DefaultLabelOptions(), "Continue")),
				ui.Button(func() { exitBtn <- true }, ui.Label(ui.DefaultLabelOptions(), "Exit to main menu")),
			)),
		))),
	}
}

func (p *pauseMenu) Draw(screen *ebiten.Image) error {
	// Dim the background with translucent black texture
	p.opts.GeoM.Reset()
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	// Fixes issue #3
	p.opts.GeoM.Scale(float64(w)+2, float64(h)+2)
	screen.DrawImage(p.tex, p.opts)

	err := p.view.Draw(screen, 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func (p *pauseMenu) ButtonPressed() buttonPressedEvent {
	if err := p.view.Update(); err != nil {
		log.Panicf("pauseMenu.ButtonPressed() - %v", err)
	}

	select {
	case <-p.continueBtn:
		log.Println("pauseMenu - \"Continue\" button pressed")
		return continueButtonPressed
	case <-p.exitBtn:
		log.Println("pauseMenu - \"Exit to main menu\" button pressed")
		return exitButtonPressed
	default:
		return noButtonPressed
	}
}
