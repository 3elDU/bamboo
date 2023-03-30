/*
	Buttons, input fields, etc.
*/

package ui

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"unicode/utf8"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/font"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// screen is the root component
type screenComponent struct {
	screen *ebiten.Image
	child  View
	id     uint64
}

func Screen(child View) *screenComponent {
	id := rand.Uint64()
	s := &screenComponent{child: child, id: id}
	child.SetParent(s)
	return s
}
func (s screenComponent) ID() uint64 {
	return s.id
}
func (s screenComponent) SetParent(parent View) {
	log.Fatal("UI - Attempted to set parent for screen component")
}
func (s *screenComponent) MaxSize() (float64, float64) {
	if s.screen == nil {
		return 0, 0
	} else {
		w, h := s.screen.Size()
		return float64(w), float64(h)
	}
}
func (s *screenComponent) CapacityForChild(child View) (float64, float64) {
	return s.MaxSize()
}
func (s *screenComponent) ComputedSize() (float64, float64) {
	return s.MaxSize()
}
func (s *screenComponent) Children() []View {
	return []View{s.child}
}
func (s *screenComponent) Update() error {
	return s.child.Update()
}
func (s *screenComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// yes, a little hacky
	s.screen = screen
	return s.child.Draw(screen, x, y)
}

type paddingComponent struct {
	baseView
	child   View
	padding float64
}

func Padding(amount float64, child View) *paddingComponent {
	p := &paddingComponent{child: child, padding: amount, baseView: newBaseView()}
	child.SetParent(p)
	return p
}
func (p *paddingComponent) MaxSize() (float64, float64) {
	return p.parent.CapacityForChild(p)
}
func (p *paddingComponent) ComputedSize() (float64, float64) {
	w, h := p.MaxSize()
	return w - p.padding*Em*2, h - p.padding*Em*2
}
func (p *paddingComponent) CapacityForChild(_ View) (float64, float64) {
	w, h := p.MaxSize()
	return w - p.padding*Em*2, h - p.padding*Em*2
}
func (p *paddingComponent) Children() []View {
	return []View{p.child}
}
func (p *paddingComponent) Update() error {
	return p.child.Update()
}
func (p *paddingComponent) Draw(screen *ebiten.Image, x, y float64) error {
	return p.child.Draw(screen, x+p.padding*Em, y+p.padding*Em)
}

type StackDirection uint

const (
	VerticalStack StackDirection = iota
	HorizontalStack
)

type StackOptions struct {
	Direction   StackDirection
	Spacing     float64   // spacing between each child
	Proportions []float64 // how much % of the parent space will each child occupy
}

type stackComponent struct {
	baseView

	opts     StackOptions
	children []View
}

func Stack(opts StackOptions, children ...View) *stackComponent {
	s := &stackComponent{
		baseView: newBaseView(),
		opts:     opts,
		children: children,
	}
	for _, child := range children {
		child.SetParent(s)
	}
	return s
}
func (s *stackComponent) MaxSize() (float64, float64) {
	return s.parent.CapacityForChild(s)
}
func (s *stackComponent) ComputedSize() (w, h float64) {
	for i, child := range s.children {
		cw, ch := child.ComputedSize()

		if s.opts.Direction == VerticalStack {
			// if it's a vstack, then width is equal to the longest child's width
			if cw > w {
				w = cw
			}
		} else {
			// if it's a hstack, then height is equal to the longest child's height
			if ch > h {
				h = ch
			}
		}

		if s.opts.Direction == VerticalStack {
			h += ch
			if i < len(s.children)-1 {
				h += s.opts.Spacing * Em
			}
		} else {
			w += ch
			if i < len(s.children)-1 {
				w += s.opts.Spacing * Em
			}
		}
	}

	return
}
func (s *stackComponent) CapacityForChild(child View) (float64, float64) {
	w, h := s.parent.CapacityForChild(child)

	// index of the child
	i := -1
	for j, c := range s.children {
		if child.ID() == c.ID() {
			i = j
			break
		}
	}
	if i == -1 {
		return 0, 0
	}

	space := 1.0
	// compute proportions
	if len(s.opts.Proportions) > 0 {
		for j, p := range s.opts.Proportions {
			// if proportion exists for given child, return it
			if j == i {
				if s.opts.Direction == VerticalStack {
					return w, h * p
				} else {
					return w * p, h
				}
			}
			space -= p
		}
	}
	// if there is no proportion for current child,
	// divide remaining space equally between all remaining children
	if s.opts.Direction == VerticalStack {
		return w, (h * space) / float64(len(s.children)-len(s.opts.Proportions))
	} else {
		return (w * space) / float64(len(s.children)-len(s.opts.Proportions)), h
	}
}

func (s *stackComponent) Children() []View {
	return s.children
}
func (s *stackComponent) AddChild(child View) {
	child.SetParent(s)
	s.children = append(s.children, child)
}
func (s *stackComponent) Update() error {
	for _, child := range s.children {
		err := child.Update()
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *stackComponent) Draw(screen *ebiten.Image, x, y float64) error {
	for _, child := range s.children {
		// there is multiple children in a stack
		// but we can't return multiple errors
		// instead, we return the first error, that we encountered
		err := child.Draw(screen, x, y)
		if err != nil {
			return err
		}

		cw, ch := child.ComputedSize()
		if s.opts.Direction == VerticalStack {
			y += ch + s.opts.Spacing*Em
		} else {
			x += cw + s.opts.Spacing*Em
		}
	}
	return nil
}

type centerComponent struct {
	baseView
	child View
}

func Center(child View) *centerComponent {
	c := &centerComponent{child: child, baseView: newBaseView()}
	child.SetParent(c)
	return c
}
func (c *centerComponent) MaxSize() (float64, float64) {
	return c.parent.CapacityForChild(c)
}
func (c *centerComponent) ComputedSize() (float64, float64) {
	return c.MaxSize()
}
func (c *centerComponent) CapacityForChild(child View) (float64, float64) {
	return c.MaxSize()
}
func (c *centerComponent) Children() []View {
	return []View{c.child}
}
func (c *centerComponent) Update() error {
	return c.child.Update()
}
func (c *centerComponent) Draw(screen *ebiten.Image, x, y float64) error {
	w, h := c.parent.CapacityForChild(c)
	// ebitenutil.DrawRect(screen, x, y, w, h, colors.Blue)
	cw, ch := c.child.ComputedSize()
	return c.child.Draw(screen, (x+x+w)/2-cw/2, (y+y+h)/2-ch/2)
}

type LabelOptions struct {
	Color color.Color
	// Font size, relative to UI scaling
	Scaling float64
}

func DefaultLabelOptions() LabelOptions {
	return LabelOptions{
		Color:   colors.Black,
		Scaling: 1,
	}
}

type labelComponent struct {
	baseView
	opts LabelOptions
	text string
}

func Label(options LabelOptions, s string) *labelComponent {
	return &labelComponent{
		text:     s,
		opts:     options,
		baseView: newBaseView(),
	}
}
func (l *labelComponent) MaxSize() (float64, float64) {
	return l.parent.CapacityForChild(l)
}
func (l *labelComponent) ComputedSize() (float64, float64) {
	w, h := font.GetStringSize(l.text, l.opts.Scaling)
	return float64(w), float64(h)
}
func (l labelComponent) CapacityForChild(_ View) (float64, float64) {
	return 0, 0
}
func (l labelComponent) Children() []View {
	return []View{}
}
func (l *labelComponent) Update() error {
	return nil
}
func (l *labelComponent) Draw(screen *ebiten.Image, x, y float64) error {
	font.RenderFontWithOptions(screen, asset_loader.DefaultFont(), l.text, x, y, l.opts.Color, l.opts.Scaling)
	return nil
}

type buttonComponent struct {
	baseView

	tex       *ebiten.Image
	tex_hover *ebiten.Image

	child View

	// transmits the button press event from Draw() to Update() (where the handler is called)
	// cleared on every Update()
	pressed bool
	handler func()
}

func Button(handler func(), child View) *buttonComponent {
	b := &buttonComponent{
		tex:       asset_loader.Texture("button").Texture(),
		tex_hover: asset_loader.Texture("button-hover").Texture(),

		child:   child,
		handler: handler,
	}
	child.SetParent(b)
	return b
}
func (b *buttonComponent) MaxSize() (float64, float64) {
	return b.ComputedSize()
}
func (b *buttonComponent) ComputedSize() (float64, float64) {
	w, h := b.tex.Size()
	return float64(w), float64(h)
}
func (b *buttonComponent) CapacityForChild(_ View) (float64, float64) {
	w, h := b.MaxSize()
	return w - Em*2, h - Em*2
}
func (b *buttonComponent) Children() []View {
	return []View{b.child}
}
func (b *buttonComponent) Update() error {
	if b.pressed {
		go b.handler()
	}
	b.pressed = false
	return nil
}
func (b *buttonComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// TODO: implement scaling
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)

	w, h := b.ComputedSize()
	cx, cy := ebiten.CursorPosition()

	// check if the cursor is hovering over the button
	mouseOver := float64(cx) > x && float64(cy) > y && float64(cx) < x+w && float64(cy) < y+h
	if mouseOver {
		screen.DrawImage(b.tex_hover, opts)
	} else {
		screen.DrawImage(b.tex, opts)
	}

	// check if button is pressed
	if mouseOver && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.pressed = true
	}

	return b.child.Draw(screen, x+Em, y+Em)
}
func (b *buttonComponent) IsPressed() bool {
	return b.pressed
}
func (b *buttonComponent) Press() {
	b.pressed = true
}

type BackgroundImageRenderingMode uint

const (
	BackgroundStretch BackgroundImageRenderingMode = iota
	BackgroundTile
)

// TODO: implement scaling
type backgroundImageComponent struct {
	baseView
	child View

	tex  *ebiten.Image
	mode BackgroundImageRenderingMode
	opts *ebiten.DrawImageOptions
}

func BackgroundImage(mode BackgroundImageRenderingMode, texture *ebiten.Image, child View) *backgroundImageComponent {
	bg := &backgroundImageComponent{
		baseView: newBaseView(),
		child:    child,
		tex:      texture,
		mode:     mode,
		opts:     &ebiten.DrawImageOptions{},
	}
	child.SetParent(bg)
	return bg
}

func (b *backgroundImageComponent) MaxSize() (float64, float64) {
	return b.parent.CapacityForChild(b)
}
func (b *backgroundImageComponent) ComputedSize() (float64, float64) {
	return b.MaxSize()
}
func (b *backgroundImageComponent) CapacityForChild(child View) (float64, float64) {
	return b.MaxSize()
}
func (b *backgroundImageComponent) Update() error {
	return b.child.Update()
}
func (b *backgroundImageComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// draw background
	b.opts.GeoM.Reset()

	switch b.mode {
	case BackgroundStretch:
		w, h := b.tex.Size()
		sw, sh := b.parent.CapacityForChild(b)
		// scale the background, so that it matches the screen size
		b.opts.GeoM.Scale(
			float64(sw)/float64(w),
			float64(sh)/float64(h),
		)
		b.opts.GeoM.Translate(x, y)
		screen.DrawImage(b.tex, b.opts)

	case BackgroundTile:
		w, h := b.tex.Size()
		sw, sh := b.parent.CapacityForChild(b)
		for tx := 0; tx < int(sw); tx += w {
			for ty := 0; ty < int(sh); ty += h {
				b.opts.GeoM.Reset()
				b.opts.CompositeMode = 0
				// TODO: make corner rendering better
				switch {
				// handle corners properly
				case int(sw)-tx < w && int(sh)-ty < h:
					b.opts.GeoM.Scale(-1, -1)
					b.opts.GeoM.Translate(x+float64(tx)+(sw-float64(tx)), y+float64(ty)+(sh-float64(ty)))
					// b.bgOpts.CompositeMode = ebiten.CompositeModeMultiply

				case int(sw)-tx < w:
					b.opts.GeoM.Scale(-1, 1)
					b.opts.GeoM.Translate(x+float64(tx)+(sw-float64(tx)), y+float64(ty))
					// b.bgOpts.CompositeMode = ebiten.CompositeModeMultiply

				case int(sh)-ty < h:
					b.opts.GeoM.Scale(1, -1)
					b.opts.GeoM.Translate(x+float64(tx), y+float64(ty)+(sh-float64(ty)))
					// b.bgOpts.CompositeMode = ebiten.CompositeModeMultiply

				default:
					b.opts.GeoM.Translate(x+float64(tx), y+float64(ty))
				}
				screen.DrawImage(b.tex, b.opts)
			}
		}
	}

	// draw child
	return b.child.Draw(screen, x, y)
}
func (b *backgroundImageComponent) Children() []View {
	return []View{b.child}
}

type backgroundColorComponent struct {
	baseView
	child View

	clr  color.Color
	tex  *ebiten.Image
	opts *ebiten.DrawImageOptions
}

func BackgroundColor(clr color.Color, child View) *backgroundColorComponent {
	// It may seem strange, that we create an entire texture, then resize it
	// Just to fill the rectangle with color
	// But documentation says, ebitenutil.DrawRect() should be used ONLY for debugging and prototyping
	// And, as of version 2.5, it is deprecated!
	// So, I guess, this is a little workaround
	tex := ebiten.NewImage(1, 1)
	tex.Fill(clr)

	bg := &backgroundColorComponent{
		baseView: newBaseView(),
		child:    child,

		clr:  clr,
		tex:  tex,
		opts: &ebiten.DrawImageOptions{},
	}
	child.SetParent(bg)

	return bg
}

func (b *backgroundColorComponent) MaxSize() (float64, float64) {
	return b.parent.CapacityForChild(b)
}
func (b *backgroundColorComponent) ComputedSize() (float64, float64) {
	return b.MaxSize()
}
func (b *backgroundColorComponent) CapacityForChild(child View) (float64, float64) {
	return b.MaxSize()
}
func (b *backgroundColorComponent) Update() error {
	return b.child.Update()
}
func (b *backgroundColorComponent) Draw(screen *ebiten.Image, x, y float64) error {
	b.opts.GeoM.Reset()
	w, h := b.ComputedSize()
	b.opts.GeoM.Scale(w, h)
	b.opts.GeoM.Translate(x, y)
	screen.DrawImage(b.tex, b.opts)

	b.child.Draw(screen, x, y)
	return nil
}
func (b *backgroundColorComponent) Children() []View {
	return []View{b.child}
}

// TODO: implement focus handling
// So that multiple input widgets at once would be possible
type inputComponent struct {
	baseView
	baseFocusView

	tex        *ebiten.Image
	texFocused *ebiten.Image
	opts       *ebiten.DrawImageOptions

	label *labelComponent

	enterKey ebiten.Key
	input    string
	handler  func(string)

	pressedKeys []rune
}

func Input(handler func(string), enterKey ebiten.Key, initialFocus bool) *inputComponent {
	return &inputComponent{
		baseView:      newBaseView(),
		baseFocusView: baseFocusView{focused: initialFocus},

		tex:        asset_loader.Texture("inputfield").Texture(),
		texFocused: asset_loader.Texture("inputfield-focused").Texture(),
		opts:       &ebiten.DrawImageOptions{},

		label: Label(DefaultLabelOptions(), ""),

		enterKey: enterKey,
		input:    "",
		handler:  handler,

		pressedKeys: make([]rune, 128),
	}
}

func (i *inputComponent) MaxSize() (float64, float64) {
	return i.ComputedSize()
}
func (i *inputComponent) ComputedSize() (float64, float64) {
	w, h := i.tex.Size()
	return float64(w), float64(h)
}
func (i *inputComponent) CapacityForChild(child View) (float64, float64) {
	w, h := i.ComputedSize()
	return w - Em*2, h - Em*2
}
func (i *inputComponent) Update() error {
	// if the element isn't focused, skip
	if !i.baseFocusView.focused {
		return nil
	}

	// listen for the key presses
	switch {
	// handle enter key
	case inpututil.IsKeyJustPressed(i.enterKey):
		go i.handler(i.input)
		i.input = ""

	// handle backspace key
	case (inpututil.IsKeyJustPressed(ebiten.KeyBackspace) ||
		inpututil.KeyPressDuration(ebiten.KeyBackspace) > 30) &&
		utf8.RuneCountInString(i.input) > 0:
		i.input = i.input[:len(i.input)-1]

	default:
		// if the input string doesn't fit in the texture, don't accept any more keys
		texCapacity, _ := i.CapacityForChild(nil)
		textSize := font.GetStringWidth(i.input, i.label.opts.Scaling)
		if textSize > int(texCapacity) {
			break
		}

		// check for pressed keys, and append them to the input string
		i.pressedKeys = ebiten.AppendInputChars(i.pressedKeys[:0])
		for _, char := range i.pressedKeys {
			i.input = fmt.Sprintf("%s%c", i.input, char) // Somewhat hacky, but works
		}
	}

	return nil
}
func (i *inputComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// check for mouse presses, and update focus accordingly
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx_int, cy_int := ebiten.CursorPosition()
		cx, cy := float64(cx_int), float64(cy_int)
		w, h := i.ComputedSize()

		if cx > x && cx < x+w && cy > y && cy < y+h {
			// if this input was pressed, set the focus to true
			i.SetFocused(true)
		} else {
			i.SetFocused(false)
		}
	}

	i.opts.GeoM.Reset()
	i.opts.GeoM.Translate(x, y)
	if i.baseFocusView.focused {
		screen.DrawImage(i.texFocused, i.opts)
	} else {
		screen.DrawImage(i.tex, i.opts)
	}

	i.label.text = i.input
	i.label.Draw(screen, x+Em, y+Em)

	return nil
}
func (i *inputComponent) Children() []View {
	return []View{}
}
func (i *inputComponent) Input() string {
	return i.input
}
func (i *inputComponent) SetInput(input string) {
	i.input = input
}
