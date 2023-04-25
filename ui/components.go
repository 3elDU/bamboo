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
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ScreenComponent is the root component
type ScreenComponent struct {
	screen *ebiten.Image
	child  View
	id     uint64
}

func Screen(child View) *ScreenComponent {
	id := rand.Uint64()
	s := &ScreenComponent{child: child, id: id}
	child.SetParent(s)
	return s
}
func (s *ScreenComponent) ID() uint64 {
	return s.id
}
func (s *ScreenComponent) SetParent(_ View) {
	log.Panicln("Attempted to set parent for screen component")
}
func (s *ScreenComponent) MaxSize() (float64, float64) {
	if s.screen == nil {
		return 0, 0
	} else {
		bounds := s.screen.Bounds()
		return float64(bounds.Dx()), float64(bounds.Dy())
	}
}
func (s *ScreenComponent) CapacityForChild(_ View) (float64, float64) {
	return s.MaxSize()
}
func (s *ScreenComponent) ComputedSize() (float64, float64) {
	return s.MaxSize()
}
func (s *ScreenComponent) Children() []View {
	return []View{s.child}
}
func (s *ScreenComponent) Update() error {
	return s.child.Update()
}
func (s *ScreenComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// yes, a little hacky
	s.screen = screen
	return s.child.Draw(screen, x, y)
}

type PaddingComponent struct {
	baseView
	child   View
	padding float64
}

func Padding(amount float64, child View) *PaddingComponent {
	p := &PaddingComponent{child: child, padding: amount, baseView: newBaseView()}
	child.SetParent(p)
	return p
}
func (p *PaddingComponent) MaxSize() (float64, float64) {
	return p.parent.CapacityForChild(p)
}
func (p *PaddingComponent) ComputedSize() (float64, float64) {
	w, h := p.MaxSize()
	return w - p.padding*Em*2, h - p.padding*Em*2
}
func (p *PaddingComponent) CapacityForChild(_ View) (float64, float64) {
	w, h := p.MaxSize()
	return w - p.padding*Em*2, h - p.padding*Em*2
}
func (p *PaddingComponent) Children() []View {
	return []View{p.child}
}
func (p *PaddingComponent) Update() error {
	return p.child.Update()
}
func (p *PaddingComponent) Draw(screen *ebiten.Image, x, y float64) error {
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

type StackComponent struct {
	baseView

	opts     StackOptions
	children []View
}

func Stack(opts StackOptions, children ...View) *StackComponent {
	s := &StackComponent{
		baseView: newBaseView(),
		opts:     opts,
		children: children,
	}
	for _, child := range children {
		child.SetParent(s)
	}
	return s
}
func (s *StackComponent) MaxSize() (float64, float64) {
	return s.parent.CapacityForChild(s)
}
func (s *StackComponent) ComputedSize() (w, h float64) {
	for i, child := range s.children {
		cw, ch := child.ComputedSize()

		switch s.opts.Direction {
		case VerticalStack:
			// vertical stack's width is equal to the widest child
			if cw > w {
				w = cw
			}
			h += ch
			// add spacing
			if i < len(s.children)-1 {
				h += s.opts.Spacing * Em
			}
		case HorizontalStack:
			// horizontal stack's height is equal to the highest child
			if ch > h {
				h = ch
			}
			w += cw
			// add spacing
			if i < len(s.children)-1 {
				w += s.opts.Spacing * Em
			}
		}
	}

	return
}
func (s *StackComponent) CapacityForChild(child View) (float64, float64) {
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

func (s *StackComponent) Children() []View {
	return s.children
}
func (s *StackComponent) AddChild(child View) {
	child.SetParent(s)
	s.children = append(s.children, child)
}
func (s *StackComponent) Update() error {
	for _, child := range s.children {
		err := child.Update()
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *StackComponent) Draw(screen *ebiten.Image, x, y float64) error {
	for _, child := range s.children {
		// there is multiple children in a stack,
		// but we can't return multiple errors.
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

type CenterComponent struct {
	baseView
	child View
}

func Center(child View) *CenterComponent {
	c := &CenterComponent{child: child, baseView: newBaseView()}
	child.SetParent(c)
	return c
}
func (c *CenterComponent) MaxSize() (float64, float64) {
	return c.parent.CapacityForChild(c)
}
func (c *CenterComponent) ComputedSize() (float64, float64) {
	return c.MaxSize()
}
func (c *CenterComponent) CapacityForChild(_ View) (float64, float64) {
	return c.MaxSize()
}
func (c *CenterComponent) Children() []View {
	return []View{c.child}
}
func (c *CenterComponent) Update() error {
	return c.child.Update()
}
func (c *CenterComponent) Draw(screen *ebiten.Image, x, y float64) error {
	w, h := c.parent.CapacityForChild(c)
	cw, ch := c.child.ComputedSize()
	return c.child.Draw(screen, x+w/2-cw/2, y+h/2-ch/2)
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

type LabelComponent struct {
	baseView
	opts LabelOptions
	text string
}

func Label(options LabelOptions, s string) *LabelComponent {
	return &LabelComponent{
		text:     s,
		opts:     options,
		baseView: newBaseView(),
	}
}
func (l *LabelComponent) MaxSize() (float64, float64) {
	return l.parent.CapacityForChild(l)
}
func (l *LabelComponent) ComputedSize() (float64, float64) {
	return font.GetStringSize(l.text, l.opts.Scaling)
}
func (l *LabelComponent) CapacityForChild(_ View) (float64, float64) {
	return 0, 0
}
func (l *LabelComponent) Children() []View {
	return []View{}
}
func (l *LabelComponent) Update() error {
	return nil
}
func (l *LabelComponent) Draw(screen *ebiten.Image, x, y float64) error {
	font.RenderFontWithOptions(screen, l.text, x, y, l.opts.Color, l.opts.Scaling)
	return nil
}

type ButtonComponent struct {
	baseView

	tex      types.Texture
	texHover types.Texture

	child View

	// transmits the button press event from Draw() to Update() (where the handler is called)
	// cleared on every Update()
	pressed bool
	handler func()
}

func Button(handler func(), child View) *ButtonComponent {
	b := &ButtonComponent{
		tex:      asset_loader.Texture("button"),
		texHover: asset_loader.Texture("button-hover"),

		child:   child,
		handler: handler,
	}
	child.SetParent(b)
	return b
}
func (b *ButtonComponent) MaxSize() (float64, float64) {
	return b.ComputedSize()
}
func (b *ButtonComponent) ComputedSize() (float64, float64) {
	return b.tex.ScaledSize()
}
func (b *ButtonComponent) CapacityForChild(_ View) (float64, float64) {
	return b.ComputedSize()
}
func (b *ButtonComponent) Children() []View {
	return []View{b.child}
}
func (b *ButtonComponent) Update() error {
	if b.pressed {
		go b.handler()
	}
	b.pressed = false
	return nil
}
func (b *ButtonComponent) Draw(screen *ebiten.Image, x, y float64) error {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	opts.GeoM.Translate(x, y)

	w, h := b.ComputedSize()
	cx, cy := ebiten.CursorPosition()

	// check if the cursor is hovering over the button
	mouseOver := float64(cx) > x && float64(cy) > y && float64(cx) < x+w && float64(cy) < y+h
	if mouseOver {
		screen.DrawImage(b.texHover.Texture(), opts)
	} else {
		screen.DrawImage(b.tex.Texture(), opts)
	}

	// check if button is pressed
	if mouseOver && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.pressed = true
	}

	// draw child in the center
	cw, ch := b.child.ComputedSize()
	return b.child.Draw(screen, x+w/2-cw/2, y+h/2-ch/2)
}
func (b *ButtonComponent) IsPressed() bool {
	return b.pressed
}
func (b *ButtonComponent) Press() {
	b.pressed = true
}

type BackgroundImageRenderingMode uint

const (
	BackgroundStretch BackgroundImageRenderingMode = iota
	BackgroundTile
)

type BackgroundImageComponent struct {
	baseView
	child View

	tex  *ebiten.Image
	mode BackgroundImageRenderingMode
	opts *ebiten.DrawImageOptions
}

func BackgroundImage(mode BackgroundImageRenderingMode, texture *ebiten.Image, child View) *BackgroundImageComponent {
	bg := &BackgroundImageComponent{
		baseView: newBaseView(),
		child:    child,
		tex:      texture,
		mode:     mode,
		opts:     &ebiten.DrawImageOptions{},
	}
	child.SetParent(bg)
	return bg
}

func (b *BackgroundImageComponent) MaxSize() (float64, float64) {
	return b.parent.CapacityForChild(b)
}
func (b *BackgroundImageComponent) ComputedSize() (float64, float64) {
	return b.MaxSize()
}
func (b *BackgroundImageComponent) CapacityForChild(_ View) (float64, float64) {
	return b.MaxSize()
}
func (b *BackgroundImageComponent) Update() error {
	return b.child.Update()
}
func (b *BackgroundImageComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// draw background
	b.opts.GeoM.Reset()

	switch b.mode {
	case BackgroundStretch:
		bounds := b.tex.Bounds()
		sw, sh := b.parent.CapacityForChild(b)
		// scale the background, so that it matches the screen size
		b.opts.GeoM.Scale(
			sw/float64(bounds.Dx()),
			sh/float64(bounds.Dy()),
		)
		b.opts.GeoM.Translate(x, y)
		screen.DrawImage(b.tex, b.opts)

	case BackgroundTile:
		bounds := b.tex.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		sw, sh := b.parent.CapacityForChild(b)
		for tx := 0; tx < int(sw); tx += w {
			for ty := 0; ty < int(sh); ty += h {
				b.opts.GeoM.Reset()
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
func (b *BackgroundImageComponent) Children() []View {
	return []View{b.child}
}

type BackgroundColorComponent struct {
	baseView
	child View

	clr  color.Color
	tex  *ebiten.Image
	opts *ebiten.DrawImageOptions
}

func BackgroundColor(clr color.Color, child View) *BackgroundColorComponent {
	// It may seem strange, that we create an entire texture, then resize it,
	// just to fill the rectangle with color.
	// But documentation says, ebitenutil.DrawRect() should be used ONLY for debugging and prototyping.
	// And, as of version 2.5, it is deprecated!
	// So, I guess, this is a little workaround
	tex := ebiten.NewImage(1, 1)
	tex.Fill(clr)

	bg := &BackgroundColorComponent{
		baseView: newBaseView(),
		child:    child,

		clr:  clr,
		tex:  tex,
		opts: &ebiten.DrawImageOptions{},
	}
	child.SetParent(bg)

	return bg
}

func (b *BackgroundColorComponent) MaxSize() (float64, float64) {
	return b.parent.CapacityForChild(b)
}
func (b *BackgroundColorComponent) ComputedSize() (float64, float64) {
	return b.MaxSize()
}
func (b *BackgroundColorComponent) CapacityForChild(_ View) (float64, float64) {
	return b.MaxSize()
}
func (b *BackgroundColorComponent) Update() error {
	return b.child.Update()
}
func (b *BackgroundColorComponent) Draw(screen *ebiten.Image, x, y float64) error {
	b.opts.GeoM.Reset()
	w, h := b.ComputedSize()
	b.opts.GeoM.Scale(w, h)
	b.opts.GeoM.Translate(x, y)
	screen.DrawImage(b.tex, b.opts)

	b.child.Draw(screen, x, y)
	return nil
}
func (b *BackgroundColorComponent) Children() []View {
	return []View{b.child}
}

type InputComponent struct {
	baseView
	baseFocusView

	tex        types.Texture
	texFocused types.Texture
	opts       *ebiten.DrawImageOptions

	label *LabelComponent

	enterKey ebiten.Key
	input    string
	handler  func(string)

	pressedKeys []rune
}

func Input(handler func(string), enterKey ebiten.Key, initialFocus bool) *InputComponent {
	return &InputComponent{
		baseView:      newBaseView(),
		baseFocusView: baseFocusView{focused: initialFocus},

		tex:        asset_loader.Texture("inputfield"),
		texFocused: asset_loader.Texture("inputfield-focused"),
		opts:       &ebiten.DrawImageOptions{},

		label: Label(DefaultLabelOptions(), ""),

		enterKey: enterKey,
		input:    "",
		handler:  handler,

		pressedKeys: make([]rune, 128),
	}
}

func (i *InputComponent) MaxSize() (float64, float64) {
	return i.ComputedSize()
}
func (i *InputComponent) ComputedSize() (float64, float64) {
	return i.tex.ScaledSize()
}
func (i *InputComponent) CapacityForChild(_ View) (float64, float64) {
	return i.ComputedSize()
}
func (i *InputComponent) Update() error {
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
		if textSize > texCapacity {
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
func (i *InputComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// check for mouse presses, and update focus accordingly
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_cx, _cy := ebiten.CursorPosition()
		cx, cy := float64(_cx), float64(_cy)
		w, h := i.ComputedSize()

		if cx > x && cx < x+w && cy > y && cy < y+h {
			// if this input was pressed, set the focus to true
			i.SetFocused(true)
		} else {
			i.SetFocused(false)
		}
	}

	i.opts.GeoM.Reset()
	i.opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	i.opts.GeoM.Translate(x, y)

	if i.baseFocusView.focused {
		screen.DrawImage(i.texFocused.Texture(), i.opts)
	} else {
		screen.DrawImage(i.tex.Texture(), i.opts)
	}

	i.label.text = i.input

	w, h := i.ComputedSize()
	cw, ch := i.label.ComputedSize()
	i.label.Draw(screen, x+w/2-cw/2, y+h/2-ch/2)

	return nil
}
func (i *InputComponent) Children() []View {
	return []View{}
}
func (i *InputComponent) Input() string {
	return i.input
}
func (i *InputComponent) SetInput(input string) {
	i.input = input
}
