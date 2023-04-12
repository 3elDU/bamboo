package scenes

import (
	"bytes"
	"encoding/binary"
	"github.com/3elDU/bamboo/types"
	"hash/fnv"
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/game"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/3elDU/bamboo/world"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type NewWorldScene struct {
	view ui.View

	// form results will be received through this channel
	// first string is world name, second is world seed
	formData chan []string
}

func NewNewWorldScene() *NewWorldScene {
	formData := make(chan []string, 1)

	return &NewWorldScene{
		formData: formData,

		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(), ui.Center(
			ui.Form(
				"Create a new world",
				formData,
				ui.FormPrompt{Title: "World name"},
				ui.FormPrompt{Title: "World seed"},
			),
		))),
	}
}

func seedFromString(s string) (seed int64) {
	hash := fnv.New64a().Sum([]byte(s))
	binary.Read(bytes.NewReader(hash), binary.BigEndian, &seed)
	return
}

func (s *NewWorldScene) Update() {
	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a viev: %v", err)
	}

	select {
	case formData := <-s.formData:
		worldName, seedString := formData[0], formData[1]
		seed := seedFromString(seedString)

		w := world.NewWorld(types.Save{
			Name:     worldName,
			BaseUUID: uuid.New(),
			UUID:     uuid.New(),
			Seed:     seed,
		})
		scene_manager.QSwitch(game.NewGameScene(w, &player.Player{X: float64(config.PlayerStartX), Y: float64(config.PlayerStartY)}))
	default:
	}
}

func (*NewWorldScene) Destroy() {
	log.Println("NewWorldScene.Destroy() called")
}

func (s *NewWorldScene) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}
