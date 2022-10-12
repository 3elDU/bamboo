package widget

import (
	"fmt"
	"image/color"

	"github.com/3elDU/bamboo/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Anchor int

const (
	TopLeft Anchor = iota
	Top
	TopRight

	Left
	Center
	Right

	BottomLeft
	Bottom
	BottomRight
)

type Text struct {
	Text   string
	Face   font.Face
	Color  color.Color
	Anchor Anchor
}

/*
Widget simply returns an image, which then will be rendered onto the screen
*/
type Widget interface {
	Anchor() Anchor
	Update()
	Render() *ebiten.Image
}

/*
TextWidget is made for rendering simple text
it is much more convenient for the widget to return just the text, and RenderTextWidget() method
will do the text rendering itself.
Why? Imagine this scenario:
We want to render a single line of text.
Yes, we still could do this the 'hard' way.
But this comes at the expense of re-creating the texture every time  we want to render the widget.
Remember, we can't know the size beforehand.
Instead of that, we would simply return the desired text,
and RenderTextWidget will render this text directly onto the screen, making our lives a lot easier.
*/
type TextWidget interface {
	Anchor() Anchor
	Update()
	Render() Text
}

// iw, ih are widget width and height
// sw, sh are destination image width and height
func widgetPosition(iw, ih, ww, wh int, anchor Anchor) (int, int) {
	switch anchor {
	default:
		return 0, 0
	case Top:
		return ww/2 - iw/2, 0
	case TopRight:
		return ww - iw, 0
	case Left:
		return 0, wh/2 - ih/2
	case Center:
		return ww/2 - iw/2, wh/2 - ih/2
	case Right:
		return ww - iw, wh/2 - ih/2
	case BottomLeft:
		return 0, wh - ih
	case Bottom:
		return ww/2 - iw/2, wh - ih
	case BottomRight:
		return ww - iw, wh - ih
	}
}

func RenderWidget(screen *ebiten.Image, widget Widget) {
	ww, wh := screen.Size()
	img := widget.Render()

	iw, ih := img.Size()

	x, y := widgetPosition(iw, ih, ww, wh, widget.Anchor())

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(img, &op)
}

func RenderTextWidget(screen *ebiten.Image, widget TextWidget) {
	ww, wh := screen.Size()
	t := widget.Render()
	bounds := text.BoundString(t.Face, t.Text)

	x, y := widgetPosition(bounds.Dx(), bounds.Dy(), ww, wh, t.Anchor)

	engine.RenderFont(screen, t.Face, t.Text, x, y, t.Color)
}

// Universal container for both types of widgets, with useful methods
type WidgetContainer struct {
	Widgets     map[string]Widget
	TextWidgets map[string]TextWidget
}

func NewWidgetContainer() *WidgetContainer {
	return &WidgetContainer{
		Widgets:     make(map[string]Widget),
		TextWidgets: make(map[string]TextWidget),
	}
}

func (container *WidgetContainer) AddWidget(name string, widget Widget) {
	container.Widgets[name] = widget
}

func (container *WidgetContainer) GetWidget(name string) Widget {
	w, exists := container.Widgets[name]
	if !exists {
		panic(fmt.Sprintf("widget with name %v doesn't exist", name))
	}
	return w
}

func (container *WidgetContainer) AddTextWidget(name string, widget TextWidget) {
	container.TextWidgets[name] = widget
}

func (container *WidgetContainer) GetTextWidget(name string) TextWidget {
	w, exists := container.TextWidgets[name]
	if !exists {
		panic(fmt.Sprintf("text widget with name %v doesn't exist", name))
	}
	return w
}

func (container *WidgetContainer) Update() {
	for _, w := range container.Widgets {
		w.Update()
	}
	for _, w := range container.TextWidgets {
		w.Update()
	}
}

func (container *WidgetContainer) Render(screen *ebiten.Image) {
	for _, w := range container.Widgets {
		RenderWidget(screen, w)
	}
	for _, w := range container.TextWidgets {
		RenderTextWidget(screen, w)
	}
}
