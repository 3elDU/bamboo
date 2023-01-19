package widgets

import (
	"fmt"
	"image/color"

	"github.com/3elDU/bamboo/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type TextureWidget struct {
	Image  *ebiten.Image
	anchor widget.Anchor
}

func (w *TextureWidget) Update() {

}

func (w *TextureWidget) Anchor() widget.Anchor {
	w.anchor++
	if w.anchor > widget.BottomRight {
		w.anchor = 0
	}

	return w.anchor
}

func (w *TextureWidget) Render() *ebiten.Image {
	return w.Image
}

type PerfWidget struct {
	Color color.Color
	Face  font.Face
}

func (w *PerfWidget) Update() {

}

func (w *PerfWidget) Anchor() widget.Anchor {
	return widget.TopRight
}

func (w *PerfWidget) Render() widget.Text {
	return widget.Text{
		Text: fmt.Sprintf("TPS: %v\nFPS: %v",
			int(ebiten.ActualTPS()), int(ebiten.ActualFPS())),
		Color:  w.Color,
		Anchor: w.Anchor(),
	}
}

type SimpleTextWidget struct {
	Text  string
	Anc   widget.Anchor
	Color color.Color
}

func (w *SimpleTextWidget) Update() {

}

func (w *SimpleTextWidget) Anchor() widget.Anchor {
	return w.Anc
}

func (w *SimpleTextWidget) Render() widget.Text {
	return widget.Text{
		Text:   w.Text,
		Color:  w.Color,
		Anchor: w.Anc,
	}
}
