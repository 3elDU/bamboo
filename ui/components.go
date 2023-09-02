/*
	Buttons, input fields, etc.
*/

package ui

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"unicode/utf8"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ScreenComponent is the root component
type ScreenComponent struct {
	// keep the reference to the screen,
	// to be able to access it during Update()
	screen *ebiten.Image
	style  ComponentStyle
	child  Component
	id     uint64
}

func Screen(child Component) *ScreenComponent {
	id := _id
	_id += 1
	s := &ScreenComponent{
		child: child,
		style: defaultStyle,
		id:    id,
	}
	child.SetParent(s)
	return s
}
func (s *ScreenComponent) ID() uint64 {
	return s.id
}
func (s *ScreenComponent) SetParent(_ Component) {
	log.Panicln("Attempted to set parent for screen component")
}
func (s *ScreenComponent) Alignment() ComponentAlignment {
	return AlignNone
}
func (s *ScreenComponent) MaxSize() (float64, float64) {
	if s.screen == nil {
		return 0, 0
	} else {
		bounds := s.screen.Bounds()
		return float64(bounds.Dx()), float64(bounds.Dy())
	}
}
func (s *ScreenComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return s.MaxSize()
}
func (s *ScreenComponent) CapacityForChild(_ Component) (float64, float64) {
	return s.MaxSize()
}
func (s *ScreenComponent) ComputedSize() (float64, float64) {
	return s.MaxSize()
}
func (s *ScreenComponent) Children() []Component {
	return []Component{s.child}
}
func (s *ScreenComponent) Update() error {
	return s.child.Update()
}
func (s *ScreenComponent) Draw(screen *ebiten.Image, x, y float64) error {
	// yes, a little hacky
	s.screen = screen
	return s.child.Draw(screen, x, y)
}
func (s *ScreenComponent) HasCustomStyles() bool {
	return false
}
func (s *ScreenComponent) Style() *ComponentStyle {
	return &s.style
}

type PaddingComponent struct {
	baseComponent
	child    Component
	paddingX float64
	paddingY float64
}

func _padding(x, y float64, child Component) *PaddingComponent {
	p := &PaddingComponent{
		baseComponent: newBaseComponent(),
		child:         child,
		paddingX:      x,
		paddingY:      y,
	}
	child.SetParent(p)
	return p
}

// Adds padding on all sides
func Padding(amount float64, child Component) *PaddingComponent {
	return _padding(amount, amount, child)
}
func PaddingX(amount float64, child Component) *PaddingComponent {
	return _padding(amount, 0, child)
}
func PaddingY(amount float64, child Component) *PaddingComponent {
	return _padding(0, amount, child)
}
func PaddingXY(amountX, amountY float64, child Component) *PaddingComponent {
	return _padding(amountX, amountY, child)
}
func (p *PaddingComponent) MaxSize() (float64, float64) {
	return p.parent.MaxCapacityForChild(p)
}
func (p *PaddingComponent) ComputedSize() (float64, float64) {
	w, h := p.child.ComputedSize()
	return w + p.paddingX*Em*2, h + p.paddingY*Em*2
}
func (p *PaddingComponent) CapacityForChild(_ Component) (float64, float64) {
	w, h := p.parent.CapacityForChild(p)
	return w - p.paddingX*Em*2, h - p.paddingY*Em*2
}
func (p *PaddingComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	w, h := p.parent.MaxCapacityForChild(p)
	return w - p.paddingX*Em*2, h - p.paddingY*Em*2
}
func (p *PaddingComponent) Children() []Component {
	return []Component{p.child}
}
func (p *PaddingComponent) Update() error {
	return p.child.Update()
}
func (p *PaddingComponent) Draw(screen *ebiten.Image, x, y float64) error {
	return p.child.Draw(screen, x+p.paddingX*Em, y+p.paddingY*Em)
}

type StackDirection uint

const (
	VerticalStack StackDirection = iota
	HorizontalStack
)

type StackOptions struct {
	Direction     StackDirection
	Spacing       float64            // spacing between each child
	Proportions   []float64          // how much % of the parent space will each child occupy
	AlignChildren ComponentAlignment // explicit alignment of all childrens
}

type StackComponent struct {
	baseComponent

	opts     StackOptions
	children []Component
}

func _stack(opts StackOptions, children []Component) *StackComponent {
	s := &StackComponent{
		baseComponent: newBaseComponent(),
		opts:          opts,
		children:      children,
	}
	for _, child := range children {
		child.SetParent(s)
	}
	return s
}
func HStack(children ...Component) *StackComponent {
	return _stack(StackOptions{Direction: HorizontalStack}, children)
}
func VStack(children ...Component) *StackComponent {
	return _stack(StackOptions{Direction: VerticalStack}, children)
}

func (s *StackComponent) WithDirection(direction StackDirection) *StackComponent {
	s.opts.Direction = direction
	return s
}
func (s *StackComponent) WithSpacing(spacing float64) *StackComponent {
	s.opts.Spacing = spacing
	return s
}
func (s *StackComponent) WithProportions(proportions ...float64) *StackComponent {
	s.opts.Proportions = proportions
	return s
}
func (s *StackComponent) AlignChildren(alignment ComponentAlignment) *StackComponent {
	s.opts.AlignChildren = alignment
	return s
}
func (s *StackComponent) WithChildren(children ...Component) *StackComponent {
	s.children = children
	for _, child := range s.children {
		child.SetParent(s)
	}
	return s
}

func (s *StackComponent) MaxSize() (float64, float64) {
	return s.parent.MaxCapacityForChild(s)
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

// w and h are parent's capacity for the stack itself
func (s *StackComponent) _capacityForChild(parentCapacityW, parentCapacityH float64, child Component) (float64, float64) {
	// Find index of the child
	childIndex := -1
	for j, c := range s.children {
		if child.ID() == c.ID() {
			childIndex = j
			break
		}
	}
	// Return if the child doesn't belong to this stack
	if childIndex == -1 {
		return 0, 0
	}

	spaceRemaining := 1.0
	// compute proportions
	if len(s.opts.Proportions) > 0 {
		for j, proportion := range s.opts.Proportions {
			// If proportion is defined for the child, return it
			if j == childIndex {
				if s.opts.Direction == VerticalStack {
					return parentCapacityW, parentCapacityH * proportion
				} else {
					return parentCapacityW * proportion, parentCapacityH
				}
			}
			spaceRemaining -= proportion
		}
	}
	// if there is no proportion defined for current child,
	// divide remaining space equally between all remaining children
	if s.opts.Direction == VerticalStack {
		return parentCapacityW, (parentCapacityH * spaceRemaining) / float64(len(s.children)-len(s.opts.Proportions))
	} else {
		return (parentCapacityW * spaceRemaining) / float64(len(s.children)-len(s.opts.Proportions)), parentCapacityH
	}
}
func (s *StackComponent) CapacityForChild(child Component) (float64, float64) {
	w, h := s.parent.CapacityForChild(s)
	return s._capacityForChild(w, h, child)
}
func (s *StackComponent) MaxCapacityForChild(child Component) (float64, float64) {
	w, h := s.parent.MaxCapacityForChild(s)
	return s._capacityForChild(w, h, child)
}

func (s *StackComponent) Children() []Component {
	return s.children
}
func (s *StackComponent) AddChild(child Component) {
	child.SetParent(s)
	s.children = append(s.children, child)
}
func (s *StackComponent) ReplaceChild(oldChild Component, newChild Component) {
	// find index of the child
	childIndex := -1
	for i, child := range s.children {
		if child.ID() == oldChild.ID() {
			childIndex = i
			break
		}
	}

	if childIndex == -1 {
		return
	}

	newChild.SetParent(s)
	s.children[childIndex] = newChild
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
	w, h := s.ComputedSize()

	for _, child := range s.children {
		cw, ch := child.ComputedSize()
		// child position on the screen
		cx, cy := x, y

		childAlignment := child.Alignment()
		// alignment of all children can be set explicitly by the container
		if s.opts.AlignChildren != AlignNone {
			childAlignment = s.opts.AlignChildren
		}

		// Align the child
		switch childAlignment {
		case AlignNone, AlignStart:
			// AlignStart is the default position, no need to modify the child position
		case AlignCenter:
			if s.opts.Direction == VerticalStack {
				cx = x + w/2 - cw/2
			} else {
				cy = y + h/2 - ch/2
			}
		case AlignEnd:
			if s.opts.Direction == VerticalStack {
				cx = x + w - cw
			} else {
				cy = y + h - ch
			}
		}

		// there is multiple children in a stack,
		// but we can't return multiple errors.
		// instead, we return the first error, that we encountered
		err := child.Draw(screen, cx, cy)
		if err != nil {
			return err
		}

		if s.opts.Direction == VerticalStack {
			y += ch + s.opts.Spacing*Em
		} else {
			x += cw + s.opts.Spacing*Em
		}
	}
	return nil
}

type CenterComponent struct {
	baseComponent
	x bool
	y bool

	child Component
}

func _center(child Component, x, y bool) *CenterComponent {
	c := &CenterComponent{
		baseComponent: newBaseComponent(),
		x:             x,
		y:             y,
		child:         child,
	}
	child.SetParent(c)
	return c
}

// Center the component on both Vertical and Horizontal axis
func Center(child Component) *CenterComponent {
	return _center(child, true, true)
}
func HCenter(child Component) *CenterComponent {
	return _center(child, true, false)
}
func VCenter(child Component) *CenterComponent {
	return _center(child, false, true)
}
func (c *CenterComponent) MaxSize() (float64, float64) {
	return c.parent.MaxCapacityForChild(c)
}
func (c *CenterComponent) ComputedSize() (w, h float64) {
	w, h = c.child.ComputedSize()
	cw, ch := c.parent.CapacityForChild(c)
	if c.x {
		w = cw
	}
	if c.y {
		h = ch
	}

	return
}
func (c *CenterComponent) CapacityForChild(_ Component) (float64, float64) {
	return c.parent.CapacityForChild(c)
}
func (c *CenterComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return c.parent.MaxCapacityForChild(c)
}
func (c *CenterComponent) Children() []Component {
	return []Component{c.child}
}
func (c *CenterComponent) Update() error {
	return c.child.Update()
}
func (c *CenterComponent) Draw(screen *ebiten.Image, x, y float64) error {
	screenX, screenY := x, y
	w, h := c.parent.CapacityForChild(c)
	cw, ch := c.child.ComputedSize()

	if c.x {
		screenX = x + w/2 - cw/2
	}
	if c.y {
		screenY = y + h/2 - ch/2
	}

	return c.child.Draw(screen, screenX, screenY)
}

type LabelComponent struct {
	baseComponent
	text string
}

func CustomLabel(s string, color color.Color, scaling float64) *LabelComponent {
	label := &LabelComponent{
		text:          s,
		baseComponent: newBaseComponent(),
	}
	label.style.TextColor = color
	label.style.TextSize = scaling
	label.style.Modified = true
	return label
}
func Label(s string) *LabelComponent {
	return &LabelComponent{
		text:          s,
		baseComponent: newBaseComponent(),
	}
}
func LabelF(format string, a ...any) *LabelComponent {
	return &LabelComponent{
		baseComponent: newBaseComponent(),
		text:          fmt.Sprintf(format, a...),
	}
}
func ColoredLabel(s string, color color.Color) *LabelComponent {
	label := &LabelComponent{
		text:          s,
		baseComponent: newBaseComponent(),
	}
	label.style.TextColor = color
	label.style.Modified = true
	return label
}

func (l *LabelComponent) WithTextColor(newColor color.Color) *LabelComponent {
	l.style.TextColor = newColor
	l.style.Modified = true
	return l
}
func (l *LabelComponent) WithTextSize(size float64) *LabelComponent {
	l.style.TextSize = size
	l.style.Modified = true
	return l
}
func (l *LabelComponent) WithoutTextShadow() *LabelComponent {
	l.style.TextShadow = false
	l.style.Modified = true
	return l
}

func (l *LabelComponent) MaxSize() (float64, float64) {
	return l.parent.CapacityForChild(l)
}
func (l *LabelComponent) ComputedSize() (float64, float64) {
	return font.GetStringSize(l.text, l.Style().TextSize)
}
func (l *LabelComponent) CapacityForChild(_ Component) (float64, float64) {
	return 0, 0
}
func (l *LabelComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return 0, 0
}
func (l *LabelComponent) Children() []Component {
	return []Component{}
}
func (l *LabelComponent) Update() error {
	return nil
}
func (l *LabelComponent) Draw(screen *ebiten.Image, x, y float64) error {
	font.RenderFontWithOptions(screen, l.text, x, y, l.Style().TextColor, l.Style().TextSize, l.Style().TextShadow)
	return nil
}
func (l *LabelComponent) Text() string {
	return l.text
}
func (l *LabelComponent) SetText(text string) {
	l.text = text
}

type ButtonComponent[T any] struct {
	baseComponent

	child Component

	// transmits the button press event from Draw() to Update() (where the value to callback is sent)
	// cleared on every Update()
	mouseOver bool
	pressed   bool
	value     T
	callback  chan T
}

func Button[T any](callback chan T, value T, child Component) *ButtonComponent[T] {
	b := &ButtonComponent[T]{
		baseComponent: newBaseComponent(),
		child:         Padding(0.3, child),
		value:         value,
		callback:      callback,
	}
	child.SetParent(b)
	return b
}
func (b *ButtonComponent[T]) MaxSize() (float64, float64) {
	return b.ComputedSize()
}
func (b *ButtonComponent[T]) ComputedSize() (float64, float64) {
	cw, ch := b.child.ComputedSize()
	cw = math.Max(MinInputWidth, cw)
	return cw + 6*config.UIScaling, ch + 6*config.UIScaling
}
func (b *ButtonComponent[T]) CapacityForChild(_ Component) (float64, float64) {
	w, h := b.parent.CapacityForChild(b)
	return w - 6*config.UIScaling, h - 6*config.UIScaling
}
func (b *ButtonComponent[T]) MaxCapacityForChild(_ Component) (float64, float64) {
	return b.parent.MaxCapacityForChild(b)
}
func (b *ButtonComponent[T]) Children() []Component {
	return []Component{b.child}
}
func (b *ButtonComponent[T]) Update() error {
	b.pressed = b.mouseOver && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
	if b.pressed {
		go func() {
			b.callback <- b.value
		}()
	}
	return nil
}
func (b *ButtonComponent[T]) Draw(screen *ebiten.Image, x, y float64) error {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	opts.GeoM.Translate(x, y)

	w, h := b.ComputedSize()
	cw, ch := b.child.ComputedSize()

	// check if the cursor is hovering over the button
	cx, cy := ebiten.CursorPosition()
	b.mouseOver = float64(cx) > x && float64(cy) > y && float64(cx) < x+w && float64(cy) < y+h
	DrawButtonBackground(screen, b.mouseOver, x, y, w-6*config.UIScaling, h-6*config.UIScaling)

	// Draw the child centered horizontally and vertically
	return b.child.Draw(screen, x+w/2-cw/2, y+h/2-ch/2)
}
func (b *ButtonComponent[T]) IsPressed() bool {
	return b.pressed
}
func (b *ButtonComponent[T]) Press() {
	b.pressed = true
}

type BackgroundImageRenderingMode uint

const (
	BackgroundStretch BackgroundImageRenderingMode = iota
	BackgroundTile
)

type BackgroundImageComponent struct {
	baseComponent
	child Component

	tex  *ebiten.Image
	mode BackgroundImageRenderingMode
	opts *ebiten.DrawImageOptions
}

func BackgroundImage(mode BackgroundImageRenderingMode, texture *ebiten.Image, child Component) *BackgroundImageComponent {
	bg := &BackgroundImageComponent{
		baseComponent: newBaseComponent(),
		child:         child,
		tex:           texture,
		mode:          mode,
		opts:          &ebiten.DrawImageOptions{},
	}
	child.SetParent(bg)
	return bg
}
func TileBackgroundImage(texture types.Texture, child Component) *BackgroundImageComponent {
	return BackgroundImage(BackgroundTile, texture.Texture(), child)
}
func StretchBackgroundImage(texture types.Texture, child Component) *BackgroundImageComponent {
	return BackgroundImage(BackgroundStretch, texture.Texture(), child)
}

func (b *BackgroundImageComponent) MaxSize() (float64, float64) {
	return b.parent.CapacityForChild(b)
}
func (b *BackgroundImageComponent) ComputedSize() (float64, float64) {
	return b.MaxSize()
}
func (b *BackgroundImageComponent) CapacityForChild(_ Component) (float64, float64) {
	return b.MaxSize()
}
func (b *BackgroundImageComponent) MaxCapacityForChild(_ Component) (float64, float64) {
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
		textureWidth, textureHeight := float64(bounds.Dx()), float64(bounds.Dy())
		screenWidth, screenHeight := b.parent.CapacityForChild(b)
		for sx := 0.0; sx < screenWidth; sx += float64(textureWidth) * config.UIScaling {
			for sy := 0.0; sy < screenHeight; sy += float64(textureHeight) * config.UIScaling {
				b.opts.GeoM.Reset()
				b.opts.GeoM.Scale(config.UIScaling, config.UIScaling)
				switch {
				// handle corners properly
				case screenWidth-sx < textureWidth && screenHeight-sy < textureHeight:
					b.opts.GeoM.Scale(-1, -1)
					b.opts.GeoM.Translate(x+sx+(screenWidth-sx), y+sy+(screenHeight-sy))

				case screenWidth-sx < textureWidth:
					b.opts.GeoM.Scale(-1, 1)
					b.opts.GeoM.Translate(x+sx+(screenWidth-sx), y+sy)

				case screenHeight-sy < textureHeight:
					b.opts.GeoM.Scale(1, -1)
					b.opts.GeoM.Translate(x+sx, y+sy+(screenHeight-sy))

				default:
					b.opts.GeoM.Translate(x+sx, y+sy)
				}
				screen.DrawImage(b.tex, b.opts)
			}
		}
	}

	// draw child
	return b.child.Draw(screen, x, y)
}
func (b *BackgroundImageComponent) Children() []Component {
	return []Component{b.child}
}

type BackgroundColorComponent struct {
	baseComponent
	child Component

	clr  color.Color
	tex  *ebiten.Image
	opts *ebiten.DrawImageOptions
	// Whether to take all the available space, or only the space occupied by the child component
	greedy bool
}

func _bgcolor(clr color.Color, child Component, greedy bool) *BackgroundColorComponent {
	// It may seem strange, that we create an entire texture, then resize it,
	// just to fill the rectangle with color.
	// But documentation says, ebitenutil.DrawRect() should be used ONLY for debugging and prototyping.
	// And, as of version 2.5, it is deprecated!
	// So, I guess, this is a little workaround
	tex := ebiten.NewImage(1, 1)
	tex.Fill(clr)

	bg := &BackgroundColorComponent{
		baseComponent: newBaseComponent(),
		child:         child,

		clr:    clr,
		tex:    tex,
		opts:   &ebiten.DrawImageOptions{},
		greedy: greedy,
	}
	child.SetParent(bg)

	return bg
}

// Same as BackgroundColor, but draws a background with alpha value. Useful for modal dialogs
func BackgroundColorAlpha(clr color.Color, alpha uint8, child Component) *BackgroundColorComponent {
	r, g, b, _ := clr.RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
	clrRgba := color.RGBA{R: r8, G: g8, B: b8, A: alpha}
	return _bgcolor(clrRgba, child, true)
}
func BackgroundColor(clr color.Color, child Component) *BackgroundColorComponent {
	return _bgcolor(clr, child, true)
}

// Similar to BackgroundColor, but takes the same space occupied by the child, essentially drawing background "behind" the child component.
func Background(clr color.Color, child Component) *BackgroundColorComponent {
	return _bgcolor(clr, child, false)
}

func (b *BackgroundColorComponent) MaxSize() (float64, float64) {
	return b.parent.MaxCapacityForChild(b)
}
func (b *BackgroundColorComponent) ComputedSize() (float64, float64) {
	if b.greedy {
		return b.MaxSize()
	} else {
		return b.child.ComputedSize()
	}
}
func (b *BackgroundColorComponent) CapacityForChild(_ Component) (float64, float64) {
	return b.parent.CapacityForChild(b)
}
func (b *BackgroundColorComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return b.parent.MaxCapacityForChild(b)
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
func (b *BackgroundColorComponent) Children() []Component {
	return []Component{b.child}
}

type InputComponent struct {
	baseComponent
	baseFocusView
	mouseOver bool

	label *LabelComponent
	child Component

	enterKey       ebiten.Key
	input          string
	maxInputLength int
	handler        func(string)

	pressedKeys []rune
}

func Input(handler func(string), enterKey ebiten.Key, initialFocus bool) *InputComponent {
	label := Label("")
	inp := &InputComponent{
		baseComponent: newBaseComponent(),
		baseFocusView: baseFocusView{focused: initialFocus},

		label: label,
		child: Padding(0.3, label),

		enterKey:       enterKey,
		input:          "",
		maxInputLength: 4096,
		handler:        handler,

		pressedKeys: make([]rune, 128),
	}
	inp.child.SetParent(inp)
	return inp
}

// Sets maximum input length
func (i *InputComponent) WithMaxInputLength(length int) *InputComponent {
	i.maxInputLength = length
	return i
}

func (i *InputComponent) MaxSize() (float64, float64) {
	return i.ComputedSize()
}
func (i *InputComponent) ComputedSize() (float64, float64) {
	cw, ch := i.child.ComputedSize()
	cw = math.Max(MinInputWidth, cw)
	return cw + 6*config.UIScaling, ch + 6*config.UIScaling
}
func (i *InputComponent) CapacityForChild(_ Component) (float64, float64) {
	w, h := i.parent.CapacityForChild(i)
	return w - 6*config.UIScaling, h - 6*config.UIScaling
}
func (i *InputComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return i.parent.MaxCapacityForChild(i)
}
func (i *InputComponent) Update() error {
	// check for mouse presses, and update focus accordingly
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		i.SetFocused(i.mouseOver)
	}

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
		// if the input is longer than maxInputLength, don't accept any more input
		if utf8.RuneCountInString(i.input) > i.maxInputLength {
			break
		}

		// check for pressed keys, and append them to the input string
		i.pressedKeys = ebiten.AppendInputChars(i.pressedKeys[:0])
		for _, char := range i.pressedKeys {
			i.input = fmt.Sprintf("%s%c", i.input, char) // Somewhat hacky, but works
		}
	}

	i.label.SetText(i.input)
	return nil
}
func (i *InputComponent) Draw(screen *ebiten.Image, x, y float64) error {
	w, h := i.ComputedSize()

	// Check if the mouse is hovering over the component
	_cx, _cy := ebiten.CursorPosition()
	cx, cy := float64(_cx), float64(_cy)

	if cx > x && cx < x+w && cy > y && cy < y+h {
		i.mouseOver = true
	} else {
		i.mouseOver = false
	}

	cw, ch := i.child.ComputedSize()

	DrawInputBackground(screen, i.baseFocusView.focused, x, y, w-6*config.UIScaling, h-6*config.UIScaling)
	i.child.Draw(screen, x+w/2-cw/2, y+h/2-ch/2)

	return nil
}
func (i *InputComponent) Children() []Component {
	return []Component{i.child}
}
func (i *InputComponent) Input() string {
	return i.input
}
func (i *InputComponent) SetInput(input string) {
	i.input = input
}

type TooltipComponent struct {
	baseComponent
	neutral bool
	child   Component
}

func Tooltip(child Component) *TooltipComponent {
	tooltip := &TooltipComponent{
		baseComponent: newBaseComponent(),
		child:         child,
	}
	child.SetParent(tooltip)
	return tooltip
}

// Uses an alternative tooltip texture made from neutral gray colors
func (tooltip *TooltipComponent) WithNeutralColor() *TooltipComponent {
	tooltip.neutral = true
	return tooltip
}

func (tooltip *TooltipComponent) MaxSize() (float64, float64) {
	return tooltip.parent.MaxCapacityForChild(tooltip)
}
func (tooltip *TooltipComponent) ComputedSize() (float64, float64) {
	cw, ch := tooltip.child.ComputedSize()
	return cw + 6*config.UIScaling, ch + 6*config.UIScaling
}
func (tooltip *TooltipComponent) CapacityForChild(_ Component) (float64, float64) {
	w, h := tooltip.parent.CapacityForChild(tooltip)
	return w - 6*config.UIScaling, h - 6*config.UIScaling
}
func (tooltip *TooltipComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	w, h := tooltip.parent.MaxCapacityForChild(tooltip)
	return w - 6*config.UIScaling, h - 6*config.UIScaling
}
func (tooltip *TooltipComponent) Children() []Component {
	return []Component{tooltip.child}
}

func (tooltip *TooltipComponent) Update() error {
	return tooltip.child.Update()
}
func (tooltip *TooltipComponent) Draw(screen *ebiten.Image, x, y float64) error {
	w, h := tooltip.child.ComputedSize()
	if tooltip.neutral {
		DrawNeutralTooltip(screen, x, y, w, h)
	} else {
		DrawTooltipBackground(screen, x, y, w, h)
	}
	return tooltip.child.Draw(screen, x+3*config.UIScaling, y+3*config.UIScaling)
}

type OverlayComponent struct {
	baseComponent
	children []Component
}

func Overlay(children ...Component) *OverlayComponent {
	overlay := &OverlayComponent{
		baseComponent: newBaseComponent(),
		children:      children,
	}
	for _, child := range children {
		child.SetParent(overlay)
	}
	return overlay
}

func (overlay *OverlayComponent) MaxSize() (float64, float64) {
	return overlay.parent.MaxCapacityForChild(overlay)
}
func (overlay *OverlayComponent) ComputedSize() (float64, float64) {
	var maxW, maxH = 0.0, 0.0

	for _, child := range overlay.children {
		w, h := child.ComputedSize()
		if w > maxW {
			maxW = w
		}
		if h > maxH {
			maxH = h
		}
	}

	return maxW, maxH
}
func (overlay *OverlayComponent) CapacityForChild(_ Component) (float64, float64) {
	return overlay.parent.CapacityForChild(overlay)
}
func (overlay *OverlayComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return overlay.parent.MaxCapacityForChild(overlay)
}
func (overlay *OverlayComponent) Children() []Component {
	return overlay.children
}
func (overlay *OverlayComponent) Update() error {
	for _, child := range overlay.children {
		if err := child.Update(); err != nil {
			return err
		}
	}
	return nil
}
func (overlay *OverlayComponent) Draw(screen *ebiten.Image, x, y float64) error {
	for _, child := range overlay.children {
		if err := child.Draw(screen, x, y); err != nil {
			return err
		}
	}
	return nil
}

type Position int

const (
	PositionTopLeft Position = iota
	PositionTop
	PositionTopRight
	PositionLeft
	PositionCenter
	PositionRight
	PositionBottomLeft
	PositionBottom
	PositionBottomRight
)

type PositionComponent struct {
	baseComponent

	position Position
	child    Component
}

// Position the component absolutely within the container
func PositionSelf(position Position, child Component) *PositionComponent {
	component := &PositionComponent{
		baseComponent: newBaseComponent(),
		position:      position,
		child:         child,
	}
	child.SetParent(component)
	return component
}

func (position *PositionComponent) MaxSize() (float64, float64) {
	return position.parent.MaxCapacityForChild(position)
}
func (position *PositionComponent) ComputedSize() (float64, float64) {
	return position.child.ComputedSize()
}
func (position *PositionComponent) CapacityForChild(_ Component) (float64, float64) {
	return position.parent.CapacityForChild(position)
}
func (position *PositionComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return position.parent.MaxCapacityForChild(position)
}
func (position *PositionComponent) Children() []Component {
	return []Component{position.child}
}
func (position *PositionComponent) Update() error {
	return position.child.Update()
}
func (position *PositionComponent) Draw(screen *ebiten.Image, x, y float64) error {
	cw, ch := position.child.ComputedSize()
	w, h := position.parent.ComputedSize()

	switch position.position {
	case PositionTopLeft:
		return position.child.Draw(screen, x, y)
	case PositionTop:
		return position.child.Draw(screen, x+w/2-cw/2, y)
	case PositionTopRight:
		return position.child.Draw(screen, x+w-cw, y)

	case PositionLeft:
		return position.child.Draw(screen, x, y+h/2-ch/2)
	case PositionCenter:
		return position.child.Draw(screen, x+w/2-cw/2, y+h/2-ch/2)
	case PositionRight:
		return position.child.Draw(screen, x+w-cw, y+h/2-ch/2)

	case PositionBottomLeft:
		return position.child.Draw(screen, x, y+h-ch)
	case PositionBottom:
		return position.child.Draw(screen, x+w/2-cw/2, y+h-ch)
	case PositionBottomRight:
		return position.child.Draw(screen, x+w-cw, y+h-ch)
	}

	return nil
}

type StyledComponent struct {
	baseComponent
	child Component
}

// Apply custom styles to a child component
func Styled(child Component) *StyledComponent {
	styled := &StyledComponent{
		baseComponent: newBaseComponent(),
		child:         child,
	}
	child.SetParent(styled)
	return styled
}

func (styled *StyledComponent) WithChild(child Component) *StyledComponent {
	child.SetParent(styled)
	styled.child = child
	return styled
}
func (styled *StyledComponent) WithTextColor(color color.Color) *StyledComponent {
	styled.style.TextColor = color
	styled.style.Modified = true
	return styled
}
func (styled *StyledComponent) WithTextSize(size float64) *StyledComponent {
	styled.style.TextSize = size
	styled.style.Modified = true
	return styled
}
func (styled *StyledComponent) WithTextShadow(shadow bool) *StyledComponent {
	styled.style.TextShadow = shadow
	styled.style.Modified = true
	return styled
}

func (styled *StyledComponent) MaxSize() (float64, float64) {
	return styled.parent.MaxCapacityForChild(styled)
}
func (styled *StyledComponent) ComputedSize() (float64, float64) {
	return styled.child.ComputedSize()
}
func (styled *StyledComponent) CapacityForChild(_ Component) (float64, float64) {
	return styled.parent.CapacityForChild(styled)
}
func (styled *StyledComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return styled.parent.MaxCapacityForChild(styled)
}
func (styled *StyledComponent) Children() []Component {
	return []Component{styled.child}
}
func (styled *StyledComponent) Update() error {
	return styled.child.Update()
}
func (styled *StyledComponent) Draw(screen *ebiten.Image, x, y float64) error {
	return styled.child.Draw(screen, x, y)
}

type ImageComponent struct {
	baseComponent

	image   *ebiten.Image
	scaling bool

	opts *ebiten.DrawImageOptions
}

// Scaling is enabled by default
func Image(image *ebiten.Image) *ImageComponent {
	return &ImageComponent{
		baseComponent: newBaseComponent(),
		image:         image,
		scaling:       true,
		opts:          &ebiten.DrawImageOptions{},
	}
}

func (image *ImageComponent) DisableScaling() *ImageComponent {
	image.scaling = false
	return image
}

func (image *ImageComponent) MaxSize() (float64, float64) {
	return image.ComputedSize()
}
func (image *ImageComponent) ComputedSize() (float64, float64) {
	w, h := image.image.Size()
	if image.scaling {
		return float64(w) * config.UIScaling, float64(h) * config.UIScaling
	} else {
		return float64(w), float64(h)
	}
}
func (image *ImageComponent) CapacityForChild(_ Component) (float64, float64) {
	return 0, 0
}
func (image *ImageComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return 0, 0
}
func (image *ImageComponent) Children() []Component {
	return []Component{}
}
func (image *ImageComponent) Update() error {
	return nil
}
func (image *ImageComponent) Draw(screen *ebiten.Image, x, y float64) error {
	image.opts.GeoM.Reset()

	if image.scaling {
		image.opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	}
	image.opts.GeoM.Translate(x, y)

	screen.DrawImage(image.image, image.opts)

	return nil
}
