/*
Various simple scenes, used in the game
*/

package scenes

import (
	"log"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/colors"
	"github.com/3elDU/bamboo/engine/scene"
	"github.com/3elDU/bamboo/engine/widget"
	"github.com/3elDU/bamboo/game/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type mainMenu struct {
	t widget.TextWidget
}

func (m *mainMenu) Update(manager *scene.SceneManager) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		manager.End()
	}
	return nil
}

func (m *mainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(colors.DarkOrange)
	widget.RenderTextWidget(screen, m.t)
}

func (m *mainMenu) Destroy() {
	log.Println("mainMenu.Destroy() called")
}

func NewMainMenuScene() *mainMenu {
	return &mainMenu{
		t: &widgets.SimpleTextWidget{
			Text:  "Press space to begin!",
			Anc:   widget.Center,
			Color: colors.Black,
			Face:  asset_loader.DefaultFont(),
		},
	}
}
