package scenes

import (
	"bytes"
	"encoding/binary"
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

type newWorldScene struct {
	view ui.View

	// form results will be received through this channel
	// first string is world name, second is world seed
	formData chan []string
}

func NewNewWorldScene() *newWorldScene {
	formData := make(chan []string, 1)

	return &newWorldScene{
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

func (s *newWorldScene) Update() {
	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a viev: %v", err)
	}

	select {
	case formData := <-s.formData:
		world_name, seed_string := formData[0], formData[1]

		// convert string to bytes -> compute hash -> convert hash to int64
		seed_bytes := []byte(seed_string)
		seed_hash_bytes := fnv.New64a().Sum(seed_bytes)
		var seed int64
		binary.Read(bytes.NewReader(seed_hash_bytes), binary.BigEndian, &seed)

		w := world.NewWorld(world_name, uuid.New(), seed)
		scene_manager.QSwitch(game.NewGameScene(w, player.Player{X: float64(config.PlayerStartX), Y: float64(config.PlayerStartY)}))
	default:
	}
}

func (*newWorldScene) Destroy() {
	log.Println("newWorldScene.Destroy() called")
}

func (s *newWorldScene) Draw(screen *ebiten.Image) {
	err := s.view.Draw(screen, 0, 0)
	if err != nil {
		log.Panicln(err)
	}
}
