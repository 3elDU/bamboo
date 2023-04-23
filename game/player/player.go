package player

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/types"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/3elDU/bamboo/config"
	"github.com/google/uuid"
)

type Player struct {
	// Note that these are block coordinates, not pixel coordinates
	X, Y                 float64
	xVelocity, yVelocity float64

	movementDirection MovementDirection
	animationFrame    uint8
	lastFrameChange   time.Time

	// Storing the selected world, so that we know what sub-world the player is currently in
	// Used to determine what sub-world to load
	SelectedWorld types.Save
}

type MovementDirection uint8

const (
	Left MovementDirection = iota
	Right
	Up
	Down
)

type MovementVector struct {
	Left, Right, Up, Down bool
}

func (mvec MovementVector) ToFloat() (vx, vy float64) {
	if mvec.Left {
		vx -= 1
	}
	if mvec.Right {
		vx += 1
	}

	if mvec.Up {
		vy -= 1
	}
	if mvec.Down {
		vy += 1
	}

	return
}

func LoadPlayer(baseUUID uuid.UUID) *Player {
	saveDir := filepath.Join(config.WorldSaveDirectory, baseUUID.String())

	f, err := os.Open(filepath.Join(saveDir, config.PlayerInfoFile))
	if err != nil {
		// if file does not exist, just return an empty object
		return &Player{X: float64(config.PlayerStartX), Y: float64(config.PlayerStartY)}
	}

	player := new(Player)
	if err := gob.NewDecoder(f).Decode(player); err != nil {
		log.Panicf("LoadPlayer() - failed to decode metadata - %v", err)
	}

	return player
}

func (player *Player) Save(metadata types.Save) {
	saveDir := filepath.Join(config.WorldSaveDirectory, metadata.BaseUUID.String())

	// storing the world player is currently in
	player.SelectedWorld = metadata

	// make a save directory, if it doesn't exist yet
	os.Mkdir(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, config.PlayerInfoFile))
	if err != nil {
		log.Panicf("failed to create player metadata file")
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(player); err != nil {
		log.Panicf("failed to write player metadata")
	}
}
