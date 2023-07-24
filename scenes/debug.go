package scenes

import (
	"image/color"
	"log"

	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DebugMenu struct {
	view ui.Component

	modeSelected chan scene_manager.Scene
}

type submenu struct {
	label   string
	tooltip string
	scene   scene_manager.Scene
}

func NewDebugMenu() *DebugMenu {
	buttonStack := ui.HStack().WithSpacing(1.0)
	modeSelected := make(chan scene_manager.Scene)

	// Store a list of all debug menus as a list rather than creating lots of buttons manually
	menus := []submenu{
		{label: "Colors", tooltip: "Displays a list of all possible colors,\nwith their complementary and analogous colors.", scene: NewColorsListDebugMenu()},
	}

	for _, menu := range menus {
		buttonStack.AddChild(ui.Button(modeSelected, menu.scene, ui.Label(menu.label)))
	}

	return &DebugMenu{
		view: ui.Screen(ui.BackgroundColor(color.White, ui.Padding(1.0,
			ui.VStack(
				ui.Label("Welcome to debug menu!"),
				buttonStack,
			),
		))),
		modeSelected: modeSelected,
	}
}

func (s *DebugMenu) Update() {
	if err := s.view.Update(); err != nil {
		log.Panicln(err)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		scene_manager.Pop()
	}

	select {
	case selected := <-s.modeSelected:
		scene_manager.PushAndSwitch(selected)
	default:
	}
}

func (s *DebugMenu) Destroy() {
	log.Println("DebugMenu.Destroy() called")
}

func (s *DebugMenu) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}
