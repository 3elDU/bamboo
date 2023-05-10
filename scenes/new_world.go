package scenes

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world_type"
	"hash/fnv"
	"log"
	"math/rand"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/game"
	"github.com/3elDU/bamboo/scene_manager"
	"github.com/3elDU/bamboo/ui"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type NewWorldScene struct {
	view ui.View

	// form results will be received through this channel
	// first string is world name, second is world seed
	formData chan []string

	goBack chan bool
}

func NewNewWorldScene() *NewWorldScene {
	formData := make(chan []string, 1)
	goBack := make(chan bool, 1)

	return &NewWorldScene{
		formData: formData,
		goBack:   goBack,

		view: ui.Screen(ui.BackgroundImage(ui.BackgroundTile, asset_loader.Texture("snow").Texture(), ui.Center(
			ui.Stack(ui.StackOptions{Direction: ui.VerticalStack},
				ui.Form(
					"Create a new world",
					formData,
					ui.FormPrompt{Title: "World name"},
					ui.FormPrompt{Title: "World seed (optional)"},
				),
				ui.Button(goBack, true, ui.Label(ui.DefaultLabelOptions(), "Go back")),
			),
		))),
	}
}

func seedFromString(s string) (seed int64) {
	if s == "" {
		// if seed string is empty, generate a random one instead
		seed = rand.Int63()
		log.Println("seed: ", seed)
	} else {
		hasher := fnv.New64a()
		hasher.Write([]byte(s))
		seed = int64(hasher.Sum64())
	}
	return
}

func (s *NewWorldScene) Update() {
	if err := s.view.Update(); err != nil {
		log.Panicf("failed to update a viev: %v", err)
	}

	select {
	case <-s.goBack:
		scene_manager.Pop()
	case formData := <-s.formData:
		worldName, seedString := formData[0], formData[1]
		seed := seedFromString(seedString)

		scene_manager.QPushAndSwitch(game.NewGameScene(types.Save{
			Name:      worldName,
			BaseUUID:  uuid.New(),
			UUID:      uuid.New(),
			Seed:      seed,
			WorldType: world_type.Overworld,
			Size:      config.SizeForWorldType(world_type.Overworld),
		}))
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
