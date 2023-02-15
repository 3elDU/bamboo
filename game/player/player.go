package player

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/config"
	"github.com/google/uuid"
)

type Player struct {
	// Note that these are block coordinates, not pixel coordinates
	X, Y                 float64
	xVelocity, yVelocity float64
	movementDirection    MovementDirection
}

type MovementDirection uint8

const (
	Still MovementDirection = iota
	Left
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

func LoadPlayer(id uuid.UUID) *Player {
	saveDir := filepath.Join(config.WorldSaveDirectory, id.String())

	f, err := os.Open(filepath.Join(saveDir, "player.gob"))
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

func (p *Player) Save(id uuid.UUID) {
	saveDir := filepath.Join(config.WorldSaveDirectory, id.String())

	// make a save directory, if it doesn't exist yet
	os.Mkdir(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, "player.gob"))
	if err != nil {
		log.Panicf("failed to create player metadata file")
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(p); err != nil {
		log.Panicf("failed to write player metadata")
	}

	log.Println("Player.Save() - saved")
}
