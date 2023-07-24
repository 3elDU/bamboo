package scenes

import (
	"fmt"
	"sort"

	"github.com/3elDU/bamboo/colors"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/maps"
)

const pixelSize = int(config.UIScaling * 32)

type ColorsListDebugMenu struct {
	sortedColorList []string
}

func NewColorsListDebugMenu() *ColorsListDebugMenu {
	// Sort the keys alphabetically so that iteration
	// through a map would be in the same order every time
	sortedColorList := maps.Keys(colors.Colors)
	sort.Strings(sortedColorList)

	return &ColorsListDebugMenu{
		sortedColorList: sortedColorList,
	}
}

func (s *ColorsListDebugMenu) colorPositionOnScreen(screenWidth int, idx int) (x int, y int) {
	x = (idx * pixelSize) % screenWidth
	y = (idx * pixelSize) / screenWidth * (pixelSize * 2)
	return
}

func (s *ColorsListDebugMenu) Draw(screen *ebiten.Image) {
	screenWidth := screen.Bounds().Dx()

	for i, key := range s.sortedColorList {
		x, y := s.colorPositionOnScreen(screenWidth, i)

		clr := colors.Colors[key]
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(pixelSize), float32(pixelSize), clr, false)
	}

	cx, cy := ebiten.CursorPosition()
	for i, key := range s.sortedColorList {
		x, y := s.colorPositionOnScreen(screenWidth, i)

		if cx >= x && cx <= x+pixelSize && cy >= y && cy <= y+pixelSize {
			clr := colors.Colors[key]
			r, g, b, _ := clr.RGBA()
			// Convert from 16-bit color to 8-bit color
			r >>= 8
			g >>= 8
			b >>= 8
			ui.DrawTextTooltip(screen, cx, cy, ui.BottomRight,
				fmt.Sprintf("R: %v, G: %v, B: %v\nColor code: %v", r, g, b, key),
			)

			break
		}
	}
}

func (s *ColorsListDebugMenu) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		scene_manager.Pop()
	}
}

func (s *ColorsListDebugMenu) Destroy() {

}
